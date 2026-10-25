package internal

import (
	"log/slog"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/trash"
	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/spferror"
)

type ModelUpdateMessage interface {
	ApplyToModel(m *model) tea.Cmd
	GetReqID() int
}

type BaseMessage struct {
	reqID int
}

func (msg BaseMessage) GetReqID() int {
	return msg.reqID
}

type PasteOperationMsg struct {
	BaseMessage

	state         processbar.ProcessState
	keepClipboard bool // the pasted items did not come from the internal clipboard
}

func NewPasteOperationMsg(state processbar.ProcessState, reqID int) PasteOperationMsg {
	return PasteOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg PasteOperationMsg) ApplyToModel(m *model) tea.Cmd {
	if !msg.keepClipboard && (msg.state == processbar.Failed || msg.state == processbar.Successful) &&
		m.clipboard.IsCut() {
		m.clipboard.Reset(false)
	}
	return nil
}

type CreateOperationMsg struct {
	BaseMessage

	state processbar.ProcessState
}

func NewCreateOperationMsg(state processbar.ProcessState, reqID int) CreateOperationMsg {
	return CreateOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg CreateOperationMsg) ApplyToModel(m *model) tea.Cmd {
	return nil
}

type DeleteOperationMsg struct {
	BaseMessage

	state processbar.ProcessState
	// items that went to the trash and can be restored with undo_delete
	trashed []trash.Result
}

func NewDeleteOperationMsg(state processbar.ProcessState, reqID int, trashed []trash.Result) DeleteOperationMsg {
	return DeleteOperationMsg{
		state:   state,
		trashed: trashed,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg DeleteOperationMsg) ApplyToModel(m *model) tea.Cmd {
	// Remove selection
	m.getFocusedFilePanel().ResetSelected()
	if len(msg.trashed) > 0 {
		m.trashUndo = append(m.trashUndo, msg.trashed)
	}
	return nil
}

type UndoDeleteMsg struct {
	BaseMessage

	restored []trash.Result
	errMsgs  []string
}

func (msg UndoDeleteMsg) ApplyToModel(m *model) tea.Cmd {
	panel := m.getFocusedFilePanel()
	for _, r := range msg.restored {
		if filepath.Dir(r.OriginalPath) == panel.Location {
			panel.TargetFile = filepath.Base(r.OriginalPath)
			break
		}
	}
	m.fileModel.UpdateFilePanelsIfNeeded(true)
	if len(msg.errMsgs) == 0 {
		return nil
	}
	content := strings.Join(msg.errMsgs, "\n") + "\n\nThese items are still in the trash."
	return func() tea.Msg {
		return NewNotifyModalMsg(notify.New(true, "Could not restore everything", content, notify.NoAction), msg.reqID)
	}
}

type ProcessBarUpdateMsg struct {
	BaseMessage

	pMsg processbar.UpdateMsg
}

func (msg ProcessBarUpdateMsg) ApplyToModel(m *model) tea.Cmd {
	cmd, err := msg.pMsg.Apply(&m.processBarModel)
	if err != nil {
		slog.Error("Error applying processbar update", "error", err)
	}
	return processCmdToTeaCmd(cmd)
}

type CompressOperationMsg struct {
	BaseMessage

	state processbar.ProcessState
}

func NewCompressOperationMsg(state processbar.ProcessState, reqID int) CompressOperationMsg {
	return CompressOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg CompressOperationMsg) ApplyToModel(_ *model) tea.Cmd {
	return nil
}

type ExtractOperationMsg struct {
	BaseMessage

	state processbar.ProcessState
}

func NewExtractOperationMsg(state processbar.ProcessState, reqID int) ExtractOperationMsg {
	return ExtractOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg ExtractOperationMsg) ApplyToModel(_ *model) tea.Cmd {
	return nil
}

type MetadataMsg struct {
	BaseMessage

	meta            metadata.Metadata
	metadataFocused bool
}

func NewMetadataMsg(meta metadata.Metadata, metadataFocused bool, reqID int) MetadataMsg {
	return MetadataMsg{
		meta:            meta,
		metadataFocused: metadataFocused,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg MetadataMsg) ApplyToModel(m *model) tea.Cmd {
	m.fileMetaData.SetMetadataCache(msg.meta, msg.metadataFocused)
	selectedItem := m.getFocusedFilePanel().GetFocusedItemPtr()
	if selectedItem == nil {
		slog.Debug("Panel empty or cursor invalid. Ignoring MetadataMsg")
		return nil
	}
	if selectedItem.Location != msg.meta.GetPath() {
		slog.Debug("MetadataMsg for older files. Ignoring",
			"currentItem", selectedItem.Location, "msgItem", msg.meta.GetPath())
		return nil
	}
	if (m.focusPanel == metadataFocus) != msg.metadataFocused {
		slog.Debug("MetadataMsg for older state. Ignoring",
			"actualFocus", m.focusPanel, "msgFocus", msg.metadataFocused)
		return nil
	}
	m.fileMetaData.SetMetadata(msg.meta, msg.metadataFocused)
	return nil
}

type NotifyModalUpdateMsg struct {
	BaseMessage

	m notify.Model
}

func NewNotifyModalMsg(m notify.Model, reqID int) NotifyModalUpdateMsg {
	return NotifyModalUpdateMsg{
		m: m,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg NotifyModalUpdateMsg) ApplyToModel(m *model) tea.Cmd {
	m.notifyModel = msg.m
	return nil
}

type SpfErrorModalUpdateMsg struct {
	BaseMessage

	m spferror.Model
}

func NewSpfErrorModalMsg(m spferror.Model, reqID int) SpfErrorModalUpdateMsg {
	return SpfErrorModalUpdateMsg{
		m: m,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg SpfErrorModalUpdateMsg) ApplyToModel(m *model) tea.Cmd {
	m.spfError = msg.m
	return nil
}
