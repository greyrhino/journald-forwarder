package kafka

import (
	"encoding/json"
	"github.com/fxfitz/journald-forwarder/journald"
	"log"
	"time"
)

func ProcessJournal(c chan journald.JournalEntry, brokers string, topic string) {
	for msg := range c {

		loggly_entry := JournalEntry{
			Pid:                     msg.Pid,
			Uid:                     msg.Uid,
			Gid:                     msg.Gid,
			Comm:                    msg.Comm,
			Exe:                     msg.Exe,
			Cmdline:                 msg.Cmdline,
			CapEffective:            msg.CapEffective,
			AuditSession:            msg.AuditSession,
			AuditLoginId:            msg.AuditLoginId,
			SystemdGroup:            msg.SystemdGroup,
			SystemdSession:          msg.SystemdSession,
			SystemdUnit:             msg.SystemdUnit,
			SystemdUserInit:         msg.SystemdUserInit,
			SystemdOwnerUid:         msg.SystemdOwnerUid,
			SystemdSlice:            msg.SystemdSlice,
			SelinuxContext:          msg.SelinuxContext,
			SourceRealtimeTimestamp: msg.SourceRealtimeTimestamp,
			BootId:                  msg.BootId,
			MachineId:               msg.MachineId,
			Hostname:                msg.Hostname,
			Transport:               msg.Transport,
			Timestamp:               time.Unix(0, msg.RealtimeTimestamp*1000),
			MonotonicTimestamp:      msg.MonotonicTimestamp,
			CoredumpUnit:            msg.CoredumpUnit,
			CoredumpUserInit:        msg.CoredumpUserInit,
			ObjectPid:               msg.ObjectPid,
			ObjectUid:               msg.ObjectUid,
			ObjectGid:               msg.ObjectGid,
			ObjectComm:              msg.ObjectComm,
			ObjectExe:               msg.ObjectExe,
			ObjectCmdline:           msg.ObjectCmdline,
			ObjectAuditSession:      msg.ObjectAuditSession,
			ObjectAuditLoginId:      msg.ObjectAuditLoginId,
			ObjectSystemdCgroup:     msg.ObjectSystemdCgroup,
			ObjectSystemdSession:    msg.ObjectSystemdSession,
			ObjectSystemdUnit:       msg.ObjectSystemdUnit,
			ObjectSystemdUserInit:   msg.ObjectSystemdUserInit,
			ObjectSystemdOwnerUid:   msg.ObjectSystemdOwnerUid,
			Message:                 msg.Message,
			MessageId:               msg.MessageId,
			Priority:                msg.Priority,
			CodeFile:                msg.CodeFile,
			CodeLine:                msg.CodeLine,
			CodeFunc:                msg.CodeFunc,
			ErrNo:                   msg.ErrNo,
			SyslogFacility:          msg.SyslogFacility,
			SyslogIdentifier:        msg.SyslogIdentifier,
		}
		json_entry, _ := json.Marshal(loggly_entry)
		// SendEvent(string(json_entry), broker, topic)

		err := LoadProducers(brokers, topic)
		if err != nil {
			log.Println("unable to load connect to kafka", err)
			return 1 // is this... right?
		}

		Send(string(json_entry))
	}
}
