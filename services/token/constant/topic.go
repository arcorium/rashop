package constant

import intev "github.com/arcorium/rashop/contract/integration/event"

const (
  TokenDomainEventTopic      = "token.dom.event"
  TokenDlqTopic              = "token.dlq"
  TokenIntegrationEventTopic = intev.TokenTopic
)

var ListeningTopics = []string{
  intev.CustomerTopic,
  intev.AuthenticationTopic,
}
