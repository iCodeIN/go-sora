// Package sora provides WebRTC signaling feature for WebRTC SFU Sora
package sora

import (
	"fmt"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v2"
)

const (
	clientVersion = "go-sora v0.2.0"
)

// DefaultOptions は Sora 接続設定のデフォルト値を生成して返します。
func DefaultOptions() *ConnectionOptions {
	return &ConnectionOptions{
		Role:        RecvOnlyRole,
		Audio:       true,
		Video:       &Video{CodecType: webrtc.VP9},
		Multistream: false,
		Debug:       false,
		Metadata:    &Metadata{},
	}
}

// CreateVideoCodec はコーデックに対応する webrtc.RTPCodec を生成して返します。
func CreateVideoCodec(codecType VideoCodecType) (*webrtc.RTPCodec, error) {
	var codec *webrtc.RTPCodec

	switch codecType {
	case VideoCodecTypeVP8:
		codec = webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000)
	case VideoCodecTypeVP9:
		codec = webrtc.NewRTPVP9Codec(webrtc.DefaultPayloadTypeVP9, 90000)
	default:
		return nil, fmt.Errorf("go-sora does not suport video codec '%s'", codecType)
	}

	return codec, nil
}

// NewConnection は Sora Connection を生成して返します。
func NewConnection(soraURL string, channelID string, options *ConnectionOptions) *Connection {
	if options == nil {
		options = DefaultOptions()
	}

	options.SoraURL = soraURL
	options.ChannelID = channelID

	c := &Connection{
		Options: options,

		connectionID: "",
		clientID:     "",

		ws:              nil,
		pc:              nil,
		pcConfig:        webrtc.Configuration{},
		connectionState: webrtc.ICEConnectionStateNew,
		answerSent:      false,

		onOpenHandler:            func(pc *webrtc.PeerConnection, m webrtc.MediaEngine) {},
		onConnectHandler:         func() {},
		onDisconnectHandler:      func(reason string, err error) {},
		onTrackHandler:           func(track *webrtc.Track) {},
		onTrackPacketHandler:     func(track *webrtc.Track, packet *rtp.Packet) {},
		onSignalingNotifyHandler: func(eventType string, message *SignalingNotifyMessage) {},
		onSpotlightNotifyHandler: func(eventType string, message *SpotlightNotifyMessage) {},
		onNetworkNotifyHandler:   func(eventType string, message *NetworkNotifyMessage) {},
		onPushHandler:            func(message []byte) {},
	}

	return c
}
