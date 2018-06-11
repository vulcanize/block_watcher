package ping

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	u "gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	host "gx/ipfs/QmQQGtcp6nVUrQjNsnU53YWV1q8fK1Kd9S7FEkYbRZzxry/go-libp2p-host"
	peer "gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
	inet "gx/ipfs/QmXdgNhVEgjLxjUoMs5ViQL7pboAt3Y7V7eGHRiE4qrmTE/go-libp2p-net"
	logging "gx/ipfs/Qmbi1CTJsbnBZjCEgc2otwu8cUFPsGpzWXG7edVCLZ7Gvk/go-log"
)

var log = logging.Logger("ping")

const PingSize = 32

const ID = "/ipfs/ping/1.0.0"

const pingTimeout = time.Second * 60

type PingService struct {
	Host host.Host
}

func NewPingService(h host.Host) *PingService {
	ps := &PingService{h}
	h.SetStreamHandler(ID, ps.PingHandler)
	return ps
}

func (p *PingService) PingHandler(s inet.Stream) {
	buf := make([]byte, PingSize)

	errCh := make(chan error, 1)
	defer close(errCh)
	timer := time.NewTimer(pingTimeout)
	defer timer.Stop()

	go func() {
		select {
		case <-timer.C:
			log.Debug("ping timeout")
			s.Reset()
		case err, ok := <-errCh:
			if ok {
				log.Debug(err)
				if err == io.EOF {
					s.Close()
				} else {
					s.Reset()
				}
			} else {
				log.Error("ping loop failed without error")
			}
		}
	}()

	for {
		_, err := io.ReadFull(s, buf)
		if err != nil {
			errCh <- err
			return
		}

		_, err = s.Write(buf)
		if err != nil {
			errCh <- err
			return
		}

		timer.Reset(pingTimeout)
	}
}

func (ps *PingService) Ping(ctx context.Context, p peer.ID) (<-chan time.Duration, error) {
	s, err := ps.Host.NewStream(ctx, p, ID)
	if err != nil {
		return nil, err
	}

	out := make(chan time.Duration)
	go func() {
		defer close(out)
		defer s.Close()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				t, err := ping(s)
				if err != nil {
					s.Reset()
					log.Debugf("ping error: %s", err)
					return
				}

				ps.Host.Peerstore().RecordLatency(p, t)
				select {
				case out <- t:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out, nil
}

func ping(s inet.Stream) (time.Duration, error) {
	buf := make([]byte, PingSize)
	u.NewTimeSeededRand().Read(buf)

	before := time.Now()
	_, err := s.Write(buf)
	if err != nil {
		return 0, err
	}

	rbuf := make([]byte, PingSize)
	_, err = io.ReadFull(s, rbuf)
	if err != nil {
		return 0, err
	}

	if !bytes.Equal(buf, rbuf) {
		return 0, errors.New("ping packet was incorrect!")
	}

	return time.Since(before), nil
}
