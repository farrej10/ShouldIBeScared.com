// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package go_movie_grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MoviemangerClient is the client API for Moviemanger service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MoviemangerClient interface {
	GetMovie(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Movie, error)
	GetMovies(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Movies, error)
}

type moviemangerClient struct {
	cc grpc.ClientConnInterface
}

func NewMoviemangerClient(cc grpc.ClientConnInterface) MoviemangerClient {
	return &moviemangerClient{cc}
}

func (c *moviemangerClient) GetMovie(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Movie, error) {
	out := new(Movie)
	err := c.cc.Invoke(ctx, "/movie.Moviemanger/GetMovie", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviemangerClient) GetMovies(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Movies, error) {
	out := new(Movies)
	err := c.cc.Invoke(ctx, "/movie.Moviemanger/GetMovies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MoviemangerServer is the server API for Moviemanger service.
// All implementations must embed UnimplementedMoviemangerServer
// for forward compatibility
type MoviemangerServer interface {
	GetMovie(context.Context, *Params) (*Movie, error)
	GetMovies(context.Context, *Params) (*Movies, error)
	mustEmbedUnimplementedMoviemangerServer()
}

// UnimplementedMoviemangerServer must be embedded to have forward compatible implementations.
type UnimplementedMoviemangerServer struct {
}

func (UnimplementedMoviemangerServer) GetMovie(context.Context, *Params) (*Movie, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovie not implemented")
}
func (UnimplementedMoviemangerServer) GetMovies(context.Context, *Params) (*Movies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovies not implemented")
}
func (UnimplementedMoviemangerServer) mustEmbedUnimplementedMoviemangerServer() {}

// UnsafeMoviemangerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MoviemangerServer will
// result in compilation errors.
type UnsafeMoviemangerServer interface {
	mustEmbedUnimplementedMoviemangerServer()
}

func RegisterMoviemangerServer(s grpc.ServiceRegistrar, srv MoviemangerServer) {
	s.RegisterService(&Moviemanger_ServiceDesc, srv)
}

func _Moviemanger_GetMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Params)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviemangerServer).GetMovie(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movie.Moviemanger/GetMovie",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviemangerServer).GetMovie(ctx, req.(*Params))
	}
	return interceptor(ctx, in, info, handler)
}

func _Moviemanger_GetMovies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Params)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviemangerServer).GetMovies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movie.Moviemanger/GetMovies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviemangerServer).GetMovies(ctx, req.(*Params))
	}
	return interceptor(ctx, in, info, handler)
}

// Moviemanger_ServiceDesc is the grpc.ServiceDesc for Moviemanger service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Moviemanger_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "movie.Moviemanger",
	HandlerType: (*MoviemangerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMovie",
			Handler:    _Moviemanger_GetMovie_Handler,
		},
		{
			MethodName: "GetMovies",
			Handler:    _Moviemanger_GetMovies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "movie/movie.proto",
}
