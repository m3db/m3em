// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package operator

import (
	"context"
	"fmt"

	"github.com/m3db/m3em/build"
	"github.com/m3db/m3em/generated/proto/m3em"
	"github.com/m3db/m3em/os/fs"

	"github.com/m3db/m3x/log"
	"google.golang.org/grpc"
)

type operator struct {
	endpoint   string
	clientConn *grpc.ClientConn
	client     m3em.OperatorClient
	opts       Options
	logger     xlog.Logger
}

// New creates a new operator
func New(
	endpoint string,
	opts Options,
) (Operator, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTimeout(opts.Timeout()), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &operator{
		endpoint:   endpoint,
		clientConn: conn,
		client:     m3em.NewOperatorClient(conn),
		opts:       opts,
		logger:     opts.InstrumentOptions().Logger(),
	}, nil
}

func (o *operator) Setup(
	bld build.ServiceBuild,
	conf build.ServiceConfiguration,
	token string,
	force bool,
) error {
	if err := o.sendSetup(token, force); err != nil {
		return fmt.Errorf("unable to reset: %v", err)
	}

	if err := o.sendBuild(bld, force); err != nil {
		return fmt.Errorf("unable to transfer build: %v", err)
	}

	if err := o.sendConfig(conf, force); err != nil {
		return fmt.Errorf("unable to transfer config: %v", err)
	}

	return nil
}

func (o *operator) sendSetup(token string, force bool) error {
	ctx := context.Background()
	_, err := o.client.Setup(ctx, &m3em.SetupRequest{
		Token: token,
		Force: force,
	})
	return err
}

func (o *operator) sendBuild(
	bld build.ServiceBuild,
	overwrite bool,
) error {
	filename := bld.SourcePath()
	iter, err := fs.NewSizedFileReaderIter(filename, o.opts.TransferBufferSize())
	if err != nil {
		return err
	}
	defer iter.Close()
	return o.sendFile(iter, m3em.FileType_M3DB_BINARY, bld.ID(), overwrite)
}

func (o *operator) sendConfig(
	conf build.ServiceConfiguration,
	overwrite bool,
) error {
	bytes, err := conf.MarshalText()
	if err != nil {
		return err
	}
	iter := fs.NewBytesReaderIter(bytes)
	return o.sendFile(iter, m3em.FileType_M3DB_CONFIG, conf.ID(), overwrite)
}

func (o *operator) sendFile(
	iter fs.FileReaderIter,
	fileType m3em.FileType,
	filename string,
	overwrite bool,
) error {
	ctx := context.Background()
	stream, err := o.client.Transfer(ctx)
	if err != nil {
		return err
	}
	chunkIdx := 0
	for ; iter.Next(); chunkIdx++ {
		bytes := iter.Current()
		request := &m3em.TransferRequest{
			Type:       fileType,
			Filename:   filename,
			Overwrite:  overwrite,
			ChunkBytes: bytes,
			ChunkIdx:   int32(chunkIdx),
		}
		err := stream.Send(request)
		if err != nil {
			stream.CloseSend()
			return err
		}
	}
	if err := iter.Err(); err != nil {
		stream.CloseSend()
		return err
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	if int(response.NumChunksRecvd) != chunkIdx {
		return fmt.Errorf("sent %d chunks, server only received %d of them", chunkIdx, response.NumChunksRecvd)
	}

	if iter.Checksum() != response.FileChecksum {
		return fmt.Errorf("expected file checksum: %d, received: %d", iter.Checksum(), response.FileChecksum)
	}

	return nil
}

func (o *operator) Start() error {
	ctx := context.Background()
	_, err := o.client.Start(ctx, &m3em.StartRequest{})
	return err
}

func (o *operator) Stop() error {
	ctx := context.Background()
	_, err := o.client.Stop(ctx, &m3em.StopRequest{})
	return err
}

func (o *operator) Teardown() error {
	ctx := context.Background()
	_, err := o.client.Teardown(ctx, &m3em.TeardownRequest{})
	return err
}

func (o *operator) Reset() error {
	panic("not implemented")
}

func (o *operator) close() error {
	if conn := o.clientConn; conn != nil {
		o.clientConn = nil
		return conn.Close()
	}
	return nil
}
