package handler

import (
    "E-Commerce/consumer-service/internal/service"
    "log"
    "os"
    "os/signal"
    "syscall"
)

type ConsumerHandler struct {
    svc service.ConsumerService
}

func NewConsumerHandler(svc service.ConsumerService) *ConsumerHandler {
    return &ConsumerHandler{svc: svc}
}

func (h *ConsumerHandler) Start() error {
    msgs, err := h.svc.ConsumeOrderCreated()
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            if err := h.svc.ProcessOrderCreated(msg); err != nil {
                log.Printf("Error processing message: %v", err)
            }
        }
    }()

    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    return nil
}