package contractus

import "testing"

func TestTransactionTyp(t *testing.T) {
	type test struct {
		name        string
		transaction Transaction
		want        string
		wantErr     bool
	}
	tests := []test{
		{
			name: "venda produtor",
			transaction: Transaction{
				Type: 1,
			},
			want: "venda produtor",
		},
		{
			name: "venda afiliado",
			transaction: Transaction{
				Type: 2,
			},
			want: "venda afiliado",
		},
		{
			name: "comissao paga",
			transaction: Transaction{
				Type: 3,
			},
			want: "comissao paga",
		},
		{
			name: "comissao recebida",
			transaction: Transaction{
				Type: 4,
			},
			want: "comissao recebida",
		},
		{
			name: "invalid transaction type",
			transaction: Transaction{
				Type: 5,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.transaction.typ()
			if (err != nil) != tt.wantErr {
				t.Errorf("typ() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got = %s; want %s", got, tt.want)
			}
		})
	}
}

func TestTransactionSellerType(t *testing.T) {
	type test struct {
		name        string
		transaction Transaction
		want        string
		wantErr     bool
	}
	tests := []test{
		{
			name: "producer",
			transaction: Transaction{
				Type: 1,
			},
			want: "producer",
		},
		{
			name: "affiliate",
			transaction: Transaction{
				Type: 2,
			},
			want: "affiliate",
		},
		{
			name: "producer",
			transaction: Transaction{
				Type: 3,
			},
			want: "producer",
		},
		{
			name: "affiliate",
			transaction: Transaction{
				Type: 4,
			},
			want: "affiliate",
		},
		{
			name: "invalid seller type",
			transaction: Transaction{
				Type: 5,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.transaction.sellerType()
			if (err != nil) != tt.wantErr {
				t.Errorf("sellerType() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got = %s; want %s", got, tt.want)
			}
		})
	}
}
