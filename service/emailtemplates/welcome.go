package emailtemplates

import (
	"allopopot-interconnect-service/models"
	"allopopot-interconnect-service/service/emailqueue"
	"fmt"
)

func GenerateWelcomeEmailTemplate(u models.User) emailqueue.EmailPayload {

	ep := new(emailqueue.EmailPayload)

	ep.To = []string{u.Email}
	ep.Subject = "Welcome to AlloPopoT Services"
	ep.Body = fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>&nbsp;</p>
		<p>Thank you for signing up to AlloPopoT Services.</p>
		<p>Please find your recovery code in the attachments of this email.</p>
		<p>&nbsp;</p>
		<p>Sincerely,</p>
		<p>AlloPopoT Identity Manager (AIM)</p>
	`, u.FirstName)

	recoverycodeAttachment := new(emailqueue.Attachments)
	recoverycodeAttachment.Filename = fmt.Sprintf("Recovery Code - %s.txt", u.Email)
	recoverycodeAttachment.MimeType = "text/plain"
	recoverycodeAttachment.SetPayload(fmt.Sprintf("This is your recovery code: %s", u.RecoveryCode))
	ep.Attachments = append(ep.Attachments, *recoverycodeAttachment)

	return *ep
}
