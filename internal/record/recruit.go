package record

type RecruitmentTarget interface {
	GetClass() string
	GetTarget() uint
}

type BasicRecruitmentTarget struct {
	Class string
	Target uint
}

func (b *BasicRecruitmentTarget) GetClass() string {return b.Class}
func (b *BasicRecruitmentTarget) GetTarget() uint {return b.Target}

func NewBasicRecruitmentTarget(r RecruitmentTarget) *BasicRecruitmentTarget {
	if r==nil {
		return nil
	} else {
		return &BasicRecruitmentTarget{
			Class:  r.GetClass(),
			Target: r.GetTarget(),
		}
	}
}