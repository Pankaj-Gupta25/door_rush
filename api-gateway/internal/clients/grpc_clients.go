package clients

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	parcelPb "github.com/sachinggsingh/PDTS-Go/pb/parcel"
	trackingPb "github.com/sachinggsingh/PDTS-Go/pb/tracking"
)

type ValidationClient struct {
	ParcelClient   parcelPb.ParcelServiceClient
	TrackingClient trackingPb.TrackingServiceClient
}

func InitGrpcClients(parcelUrl string, trackingUrl string) (*ValidationClient, error) {
	// Parcel Service Connection
	parcelConn, err := grpc.NewClient(parcelUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to Parcel Service at %s: %v", parcelUrl, err)
		return nil, err
	}
	parcelClient := parcelPb.NewParcelServiceClient(parcelConn)

	// Tracking Service Connection
	trackingConn, err := grpc.NewClient(trackingUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to Tracking Service at %s: %v", trackingUrl, err)
		// Close parcelConn if tracking fails? ideally yes, but for now simple
		return nil, err
	}
	trackingClient := trackingPb.NewTrackingServiceClient(trackingConn)

	return &ValidationClient{
		ParcelClient:   parcelClient,
		TrackingClient: trackingClient,
	}, nil
}
