package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
)

type EntityWriter struct {
	Indent          string
	FieldNameLength int

	fieldNameFormat string
	sb              *strings.Builder
}

func NewEntityWriter(indent string, fieldNameLen int) *EntityWriter {
	w := &EntityWriter{
		Indent:          indent,
		FieldNameLength: fieldNameLen,
		fieldNameFormat: fmt.Sprintf("%%-%ds", fieldNameLen),
		sb:              &strings.Builder{},
	}
	return w
}
func DefaultEntityWriter() *EntityWriter { return NewEntityWriter("", 15) }

func (w *EntityWriter) String() string { return w.sb.String() }
func (w *EntityWriter) Print()         { fmt.Print(w.String()) }
func (w *EntityWriter) Println()       { fmt.Println(w.String()) }

func (w *EntityWriter) WithIndent(value string) *EntityWriter {
	return &EntityWriter{
		Indent:          w.Indent + value,
		FieldNameLength: w.FieldNameLength,
		fieldNameFormat: w.fieldNameFormat,
		sb:              w.sb,
	}
}
func (w *EntityWriter) Clone() *EntityWriter { return w.WithIndent("") }

func (w *EntityWriter) WriteString(s string) *EntityWriter {
	w.sb.WriteString(s)
	return w
}
func (w *EntityWriter) WriteIndent() *EntityWriter { return w.WriteString(w.Indent) }
func (w *EntityWriter) Format(format string, args ...any) *EntityWriter {
	return w.WriteString(fmt.Sprintf(format, args...))
}
func (w *EntityWriter) WriteObject(objectType string, fieldWriter func(fieldWriter *EntityWriter)) *EntityWriter {
	fieldsBlock := w.WriteString(objectType).WriteString("{\n").WithIndent("  ")
	fieldWriter(fieldsBlock)
	return w.WriteIndent().WriteString("}")
}
func (w *EntityWriter) WriteArray(arrayName, itemType string, length int, arrayWriter func(index int, itemWriter *EntityWriter)) *EntityWriter {
	return w.
		WriteString(arrayName).WriteString(": ").
		WriteString("[]").
		WriteObject(itemType, func(fieldWriter *EntityWriter) {
			for i := 0; i < length; i++ {
				arrayWriter(i, fieldWriter.WriteIndent())
				fieldWriter.WriteString(",\n")
			}
		})
}
func (w *EntityWriter) WriteFieldName(fieldName string) *EntityWriter {
	return w.Format(w.fieldNameFormat, fieldName+":")
}
func (w *EntityWriter) FormatField(fieldName string, valueFormat string, args ...any) *EntityWriter {
	return w.
		WriteIndent().
		WriteFieldName(fieldName).
		Format(valueFormat, args...).
		WriteString(",\n")
}
func (w *EntityWriter) ConditionalWrite(
	condition bool,
	action, elseWriter func(conditionalWriter *EntityWriter),
) *EntityWriter {
	if condition {
		action(w)
	} else if elseWriter != nil {
		elseWriter(w)
	}
	return w
}

func (w *EntityWriter) writeModel(model *entities.Model) *EntityWriter {
	return w.
		FormatField("ID", "%q", model.ID.String()).
		FormatField("CreatedAt", "%q", model.CreatedAt.String()).
		FormatField("UpdatedAt", "%q", model.UpdatedAt).
		ConditionalWrite(
			model.DeletedAt.Valid,
			func(conditionalWriter *EntityWriter) {
				conditionalWriter.
					FormatField("DeletedAt", "%q", model.DeletedAt.Time.String())
			},
			nil,
		)
}
func (w *EntityWriter) WriteEvent(event *entities.Event) *EntityWriter {
	return w.WriteObject("&entities.Event", func(fieldWriter *EntityWriter) {
		fieldWriter.
			writeModel(&event.Model).
			FormatField("ApplicationID", "%q", event.ApplicationID.String()).
			FormatField("SessionID", "%q", event.SessionID.String()).
			FormatField("EventType", "%q", event.EventType).
			FormatField("EventData", "%q", event.EventData)
	})
}
func (w *EntityWriter) WriteSession(session *entities.Session) *EntityWriter {
	return w.WriteObject("&entities.Session", func(fieldWriter *EntityWriter) {
		fieldWriter.
			writeModel(&session.Model).
			FormatField("ApplicationID", "%q", session.ApplicationID.String()).
			FormatField("CloseTime", "%q", session.CloseTime.String()).
			FormatField("IsClosed()", "%v", session.CloseTime.Before(time.Now().UTC())).
			FormatField("IP", "%#v", session.IP).
			FormatField("OS", "%#v", session.OS).
			FormatField("Browser", "%#v", session.Browser).
			FormatField("Screen", "types.Size{Width: %v, Height: %v}", session.Screen.Width, session.Screen.Height).
			ConditionalWrite(
				len(session.Events) != 0,
				func(conditionalWriter *EntityWriter) {
					conditionalWriter.
						WriteIndent().
						WriteArray(
							"Events",
							"*entities.Event",
							len(session.Events),
							func(index int, itemWriter *EntityWriter) {
								itemWriter.WriteEvent(session.Events[index])
							}).
						WriteString(",\n")
				},
				nil,
			)
	})
}
func (w *EntityWriter) WriteApplication(application *entities.Application) *EntityWriter {
	return w.WriteObject("&entities.Application", func(fieldWriter *EntityWriter) {
		fieldWriter.
			writeModel(&application.Model).
			FormatField("Name", "%q", application.Name).
			ConditionalWrite(
				len(application.Sessions) != 0,
				func(conditionalWriter *EntityWriter) {
					conditionalWriter.
						WriteIndent().
						WriteArray(
							"Sessions",
							"*entities.Session",
							len(application.Sessions),
							func(index int, itemWriter *EntityWriter) {
								itemWriter.WriteSession(application.Sessions[index])
							}).
						WriteString(",\n")
				},
				nil).
			ConditionalWrite(
				len(application.Events) != 0,
				func(conditionalWriter *EntityWriter) {
					conditionalWriter.
						WriteIndent().
						WriteArray(
							"Events",
							"*entities.Event",
							len(application.Events),
							func(index int, itemWriter *EntityWriter) {
								itemWriter.WriteEvent(application.Events[index])
							}).
						WriteString(",\n")
				},
				nil)
	})
}
