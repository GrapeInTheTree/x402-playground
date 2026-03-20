package learn

// Topic represents a learning topic with markdown content.
type Topic struct {
	Title       string
	Description string
	Icon        string
	Content     func() string
}

// AllTopics returns the ordered list of learning topics.
func AllTopics() []Topic {
	return []Topic{
		{
			Title:       "What is x402?",
			Description: "HTTP 402 기반 머신 결제 프로토콜",
			Icon:        "🌐",
			Content:     ContentProtocol,
		},
		{
			Title:       "HTTP 402 & Payment Headers",
			Description: "결제 요구/응답 헤더 구조",
			Icon:        "📋",
			Content:     Content402,
		},
		{
			Title:       "EIP-3009: transferWithAuthorization",
			Description: "USDC 네이티브 메타 트랜잭션",
			Icon:        "✍️",
			Content:     ContentEIP3009,
		},
		{
			Title:       "Permit2: Universal ERC-20",
			Description: "Uniswap Permit2를 통한 범용 토큰 전송",
			Icon:        "🔑",
			Content:     ContentPermit2,
		},
		{
			Title:       "EIP-712: Typed Structured Data",
			Description: "구조화된 데이터 서명 표준",
			Icon:        "📐",
			Content:     ContentEIP712,
		},
		{
			Title:       "Facilitator Role & Gas Sponsorship",
			Description: "가스비 대납과 결제 중개자 역할",
			Icon:        "⛽",
			Content:     ContentFacilitator,
		},
	}
}
