package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"
)

// HandleWebhook returns a HandlerFunc that sends messages for users via a channel
func HandleWebhook(messages chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var webhook notify.WebhookMessage

		decoder := json.NewDecoder(r.Body)
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()

		if err := decoder.Decode(&webhook); err != nil {
			log.Printf("failed to decode webhook message: %v\n", err)
		}

		for _, webAlert := range webhook.Alerts {
			labels := make(map[model.LabelName]model.LabelValue)
			for k, v := range webAlert.Labels {
				labels[model.LabelName(k)] = model.LabelValue(v)
			}

			annotations := make(map[model.LabelName]model.LabelValue)
			for k, v := range webAlert.Annotations {
				annotations[model.LabelName(k)] = model.LabelValue(v)
			}

			alert := types.Alert{
				Alert: model.Alert{
					StartsAt:     webAlert.StartsAt,
					EndsAt:       webAlert.EndsAt,
					GeneratorURL: webAlert.GeneratorURL,
					Labels:       labels,
					Annotations:  annotations,
				},
			}

			var out string
			out = out + AlertMessage(alert) + "\n"

			messages <- out
		}

		w.WriteHeader(http.StatusOK)
	}
}