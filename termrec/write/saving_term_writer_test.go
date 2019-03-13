package write_test

import (
	"testing"
	"time"

	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/formatters"
	"github.com/theparanoids/ashirt/termrec/write"

	"github.com/stretchr/testify/assert"
)

func TestSaveTermWriterWriteHeader(t *testing.T) {
	writer := write.NewSaveTermWrier()

	meta := formatters.Metadata{Title: "A Header!"}
	writer.WriteHeader(meta)

	assert.Equal(t, meta, *writer.HeaderMetadata)
}

func TestSaveTermWriterWriteFooter(t *testing.T) {
	writer := write.NewSaveTermWrier()

	meta := formatters.Metadata{Title: "A Footer!"}
	writer.WriteFooter(meta)

	assert.Equal(t, meta, *writer.FooterMetadata)
}

func TestSaveTermWriterWriteHeaderAndFooter(t *testing.T) {
	writer := write.NewSaveTermWrier()

	metaH := formatters.Metadata{Title: "A Header!"}
	metaF := formatters.Metadata{Title: "A Footer!"}
	writer.WriteHeader(metaH)
	writer.WriteFooter(metaF)

	assert.Equal(t, metaH, *writer.HeaderMetadata)
	assert.Equal(t, metaF, *writer.FooterMetadata)

}

func TestSaveTermWriterWriteEvent(t *testing.T) {
	writer := write.NewSaveTermWrier()

	evt1 := common.Event{Type: "i", Data: "someData", When: 2 * time.Second}
	evt2 := common.Event{Type: "o", Data: "moreData", When: 4 * time.Second}

	writer.WriteEvent(evt1)
	assert.Equal(t, len(*writer.AllEvents), 1)
	assert.Equal(t, (*writer.AllEvents)[0], evt1)

	writer.WriteEvent(evt2)
	assert.Equal(t, len(*writer.AllEvents), 2)
	assert.Equal(t, (*writer.AllEvents)[1], evt2)
}
