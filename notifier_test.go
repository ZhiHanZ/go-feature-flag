package ffclient

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/thomaspoignant/go-feature-flag/notifier/logs"
	"github.com/thomaspoignant/go-feature-flag/notifier/slack"
	"github.com/thomaspoignant/go-feature-flag/notifier/webhook"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspoignant/go-feature-flag/notifier"

	"github.com/thomaspoignant/go-feature-flag/internal"
)

func TestGoFeatureFlag_getNotifiers(t *testing.T) {
	parsedURL, _ := url.Parse("http://webhook.com/hook")
	hostname, _ := os.Hostname()

	type fields struct {
		config Config
	}
	tests := []struct {
		name    string
		fields  fields
		want    []notifier.Notifier
		wantErr bool
	}{
		{
			name: "log + webhook notifier",
			fields: fields{
				config: Config{
					Logger: log.New(os.Stdout, "", 0),
					Notifiers: []NotifierConfig{
						&WebhookConfig{
							EndpointURL: parsedURL.String(),
							Secret:      "Secret",
							Meta: map[string]string{
								"my-app":   "go-ff-test",
								"hostname": hostname,
							},
						},
						&SlackNotifier{
							SlackWebhookURL: parsedURL.String(),
						},
					},
				},
			},
			want: []notifier.Notifier{
				&logs.Notifier{Logger: log.New(os.Stdout, "", 0)},
				&webhook.Notifier{
					Logger: log.New(os.Stdout, "", 0),
					HTTPClient: &http.Client{
						Timeout: 10 * time.Second,
					},
					EndpointURL: *parsedURL,
					Secret:      "Secret",
					Meta: map[string]string{
						"my-app":   "go-ff-test",
						"hostname": hostname,
					},
				},
				&slack.Notifier{
					Logger:     log.New(os.Stdout, "", 0),
					HTTPClient: internal.DefaultHTTPClient(),
					WebhookURL: *parsedURL,
				},
			},
		},
		{
			name: "error in DNS",
			fields: fields{
				config: Config{
					Logger: log.New(os.Stdout, "", 0),
					Notifiers: []NotifierConfig{
						&WebhookConfig{
							EndpointURL: " https://example.com/hook",
							Secret:      "Secret",
							Meta: map[string]string{
								"my-app":   "go-ff-test",
								"hostname": hostname,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNotifiers(tt.fields.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
