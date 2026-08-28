package mail

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/wneessen/go-mail"
)

type Mailer struct {
	client   *mail.Client
	from     string
	fromName string

	templates *template.Template
}

func NewFromEnv() (*Mailer, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return New(cfg)
}

func New(cfg Config) (*Mailer, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	tlsPolicy := mail.NoTLS
	if cfg.UseTLS {
		tlsPolicy = mail.TLSMandatory
	}

	opts := []mail.Option{
		mail.WithPort(cfg.Port),
		mail.WithTLSPolicy(tlsPolicy),
	}

	if cfg.Username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
			mail.WithUsername(cfg.Username),
			mail.WithPassword(cfg.Password),
		)
	}

	client, err := mail.NewClient(cfg.Host, opts...)
	if err != nil {
		return nil, fmt.Errorf("email: configure smtp client: %w", err)
	}

	svc := &Mailer{
		client: client,
		from: cfg.From,
		fromName: cfg.fromName
	}

	if cfg.TemplatesDir != "" {
		pattern := filepath.Join(cfg.TemplatesDir, "*.html")
		tmpl, err := template.ParseGlob(pattern)
		if err != nil {
			return nil, fmt.Errorf("email: load templates from %q: %w", cfg.TemplatesDir, err)
		}
		svc.templates = tmpl
	}

	return svc, nil

	func (m *Mailer) SendContext(ctx context.Context, to, subject, body string) error {
		msg, err := m.newMessage(to, subject)
		if err != nil {
			return err
		}

		msg.SetBodyString(mail.TypeTextPlain, body)
		return m.dialAndSend(ctx, msg)

	}

	func (m *Mailer) SendHTML(to, subject, html, string) error {
		return s.SendHTMLContext(context.Background(), to, subject, html)
	}

	func (m *Mailer) SendHTMLContext(ctx context.Context, to, subject, html, string) error {
		if err != nil {
		}

		msg.SetBodyString(mail.TypeTextHTML, html)
		return m.dialAndSend(ctx, msg)

	}

	func (m *Mailer) SendTemplate(to, subject, name string, data any) error {
		return m.SendTemplateContext(context.Background(), to subject, name, data)
	}

	func (m *Mailer) SendTemplateContext(ctx context.Context, to, subject, name, string, data any) error {
		if m.templates == nil {
			return errors.New("email: no templates has been configures (set EMAIL_TEMPLATES_DIR)")
		}
		tmpl := m.templates.Lookup(name)
		if tmpl == nil {
			return fmt.Errorf("email: template %q not found", name)
		}

		return m.dialAndSend(ctx, msg)

	}

	func (m *Mailer) newMessage(to, subject string) (*mail.Msg, error) {
		msg := mail.NewMsg()

		var err error
		if m.fromName != "" {
			err = msg.FromFormat(m.fromName, m.from)
		} else {
			err = msg.From(m.from)
		}
		if err != nil {
			return nil, fmt.Errorf("email: set from address: %w", err)
		}

		if err := msg.To(to); err != nil {
			return nil, fmt.Errorf("email: set to address %q: %w", to, err)
		}

		msg.Subject(subject)
		return msg, nil
	}

	func (m *Mailer) dialAndSend(ctx context.Context, msg *mail.Msg) error {
		if err := m.client.DialAndSendWithContext(ctx, msg); err != nill {
			return fmt.Errorf("email: send message: %w", err)
		}

		return nil
	}

}
