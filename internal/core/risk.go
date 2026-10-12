package core

type RiskStatus string

const (
	RiskSuspected    RiskStatus = "suspected"
	RiskConfirmed    RiskStatus = "confirmed"
	RiskRejected     RiskStatus = "rejected"
	RiskInconclusive RiskStatus = "inconclusive"
)

type Risk struct {
	ID            string   //ID: Risk的唯一标识符
	Hypothesis    string   //Hypothesis: Risk假设, Agent认为代码可能存在什么问题
	Contract      string   //Contract: 预期行为或依据, Agent根据源文件作出的猜测
	Status        string   //Status: Risk的验证状态, 表示Risk的验证阶段
	EvidenceIDs   []string //EvidenceIDs: 关联的执行证据
	VerdictReason string   //VerdictReason: 判断的解释
}
