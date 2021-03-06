/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package init

import (
	"bytes"
	"math/big"

	"github.com/ontio/layer2/node/common"
	"github.com/ontio/layer2/node/smartcontract/service/native/auth"
	params "github.com/ontio/layer2/node/smartcontract/service/native/global_params"
	"github.com/ontio/layer2/node/smartcontract/service/native/ong"
	"github.com/ontio/layer2/node/smartcontract/service/native/ont"
	"github.com/ontio/layer2/node/smartcontract/service/neovm"
	vm "github.com/ontio/layer2/node/vm/neovm"
)

func init() {
	ong.InitOng()
	ont.InitOnt()
	params.InitGlobalParams()
	auth.Init()
}

func InitBytes(addr common.Address, method string) []byte {
	bf := new(bytes.Buffer)
	builder := vm.NewParamsBuilder(bf)
	builder.EmitPushByteArray([]byte{})
	builder.EmitPushByteArray([]byte(method))
	builder.EmitPushByteArray(addr[:])
	builder.EmitPushInteger(big.NewInt(0))
	builder.Emit(vm.SYSCALL)
	builder.EmitPushByteArray([]byte(neovm.NATIVE_INVOKE_NAME))

	return builder.ToArray()
}
