// +build server

package serverrpcs

import (
	"context"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"strings"
	"time"
)

func initialCap(s string) string {
	if s=="" {return ""}
	return strings.ToUpper(s[:1])+strings.ToLower(s[1:])
}

type serverMembersHandler struct {}

func (s serverMembersHandler) GetMember(ctx context.Context, name string) (record.Member, error) {
	return dao.GetMember(initialCap(name))
}

func (s serverMembersHandler) GetMembers(ctx context.Context) (map[string]record.Member, error) {
	return dao.GetMembers()
}

func (s serverMembersHandler) MergeMember(ctx context.Context, member record.Member) (record.Member, error) {
	_, err := requiresAdmin(ctx)
	if err!=nil {
		return nil, err
	}

	if member==nil {
		return nil, nil
	}
	err = db.MakeUpdate(func(tx *bbolt.Tx) error {
		return txMergeMember(tx, member)
	})
	if err!=nil {
		return nil, err
	}
	return dao.GetMember(member.GetName())
}

func (s serverMembersHandler) MergeMembers(ctx context.Context, members []record.Member) (map[string]record.Member, error) {
	_, err := requiresAdmin(ctx)
	if err!=nil {
		return nil, err
	}

	err = db.MakeUpdate(func(tx *bbolt.Tx) error {
		// Merge twice so we get owners better
		for _, v := range members {
			if v==nil {
				continue
			}
			err := txMergeMember(tx, v)
			if err!=nil {
				return err
			}
		}

		for _, v := range members {
			if v==nil {
				continue
			}
			err := txMergeMember(tx, v)
			if err!=nil {
				return err
			}
		}

		return nil
	})
	if err!=nil {
		return nil, err
	}
	return dao.GetMembers()
}

func txMergeMember(tx *bbolt.Tx, member record.Member) error {
	if member.GetName()=="" {
		return errors.New("Each member must have a non-empty name")
	}
	if _, ok := eqspec.ClassMap[member.GetClass()]; !ok {
		return errors.New("Unrecognized class "+member.GetClass())
	}
	realOwner:=""
	_, err := dao.TxGetMember(tx, initialCap(member.GetOwner()))
	if err==nil {
		realOwner=initialCap(member.GetOwner())
	}

	oldMember, err := dao.TxGetMember(tx, initialCap(member.GetName()))
	if err==bolthold.ErrNotFound {
		return dao.TxUpsertMember(tx, &record.BasicMember{
			Name:       initialCap(member.GetName()),
			Class:      member.GetClass(),
			Level:      member.GetLevel(),
			Rank: 		member.GetRank(),
			Alt:        member.IsAlt(),
			DKP:        0,
			LastActive: time.Time{},
			Owner:      realOwner,
		})
	} else if err!=nil {
		return err
	} else {
		if realOwner=="" {
			realOwner=oldMember.GetOwner()
		}
		return dao.TxUpsertMember(tx, &record.BasicMember{
			Name:       initialCap(member.GetName()),
			Class:      member.GetClass(),
			Level:      member.GetLevel(),
			Rank:       member.GetRank(),
			Alt:        member.IsAlt(),
			DKP:        oldMember.GetDKP(),
			LastActive: oldMember.GetLastActive(),
			Owner:      realOwner,
		})
	}
}

func init() {
	register(restidl.HandleMembers(&serverMembersHandler{}))
}