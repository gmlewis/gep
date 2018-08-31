// Copyright 2014 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package boolNodes

import (
	"github.com/gmlewis/gep/v2/functions"
)

// BoolNorOnly defines the boolean functions available for this model.
var BoolNorOnly = functions.FuncMap{
	"Not":  BoolNode{10000, "Not", 1, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, x0) }},
	"And":  BoolNode{10001, "And", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x0), gepNor(x1, x1)) }},
	"Or":   BoolNode{10002, "Or", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x1), gepNor(x0, x1)) }},
	"Nand": BoolNode{10003, "Nand", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), (false)) }},
	"Nor":  BoolNode{10004, "Nor", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, x1) }},
	"Xor":  BoolNode{10005, "Xor", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x1, x1))) }},
	"Nxor": BoolNode{10006, "Nxor", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x1)), gepNor(gepNor(x0, x0), x1)) }},
	"And3": BoolNode{10007, "And3", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x2, x2))), gepNor(x1, x1))
	}},
	"Or3": BoolNode{10008, "Or3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(gepNor(x0, x1), x2), x2), (false)) }},
	"Nand3": BoolNode{10009, "Nand3", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x0), gepNor(x1, x1))), gepNor(x2, x2)), (false))
	}},
	"Nor3": BoolNode{10010, "Nor3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, gepNor(x0, gepNor(x1, x2))) }},
	"Odd3": BoolNode{10011, "Odd3", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(gepNor(x1, x1), x2)), gepNor(x0, x2)), gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x1, x2)), gepNor(x0, x2)), x1))
	}},
	"Even3": BoolNode{10012, "Even3", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x1, x2)), gepNor(x0, x2)), gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(gepNor(x1, x1), x2)), gepNor(gepNor(x0, x2), x2)))
	}},
	"And4": BoolNode{10013, "And4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x0), x2), gepNor(gepNor(x0, x0), gepNor(x3, x3))), gepNor(gepNor(x0, gepNor(x0, x3)), gepNor(gepNor(x1, x1), gepNor(x2, x2))))
	}},
	"Or4": BoolNode{10014, "Or4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(x0, x2)), gepNor(gepNor(x1, x3), x2)), (false))
	}},
	"Nand4": BoolNode{10015, "Nand4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(gepNor(gepNor(x2, x2), gepNor(x3, x3)), (false))), (false)), gepNor(x1, x1)), (false))
	}},
	"Nor4": BoolNode{10016, "Nor4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x1)), gepNor(x1, gepNor(x2, x3)))
	}},
	"Odd4": BoolNode{10017, "Odd4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, gepNor(gepNor(x0, x1), gepNor(gepNor(gepNor(x0, x2), x3), gepNor(gepNor(x0, x3), x2))))), gepNor(gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), x1), gepNor(gepNor(x2, x2), gepNor(x3, x3))), (false)), gepNor(x2, x3)), gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x1, x2))), gepNor(x1, gepNor(x1, x3))), gepNor(gepNor(gepNor(x0, x1), gepNor(x2, gepNor(x2, x3))), x1))))
	}},
	"Even4": BoolNode{10018, "Even4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, gepNor(gepNor(x0, gepNor(gepNor(x1, gepNor(x1, x2)), gepNor(gepNor(x1, x2), x2))), gepNor(x3, x3)))), gepNor(gepNor(gepNor(gepNor(x0, (false)), x3), gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x2)), gepNor(x1, x2))), gepNor(gepNor(x0, x3), gepNor(gepNor(x0, x0), gepNor(gepNor(gepNor(x1, x3), gepNor(x2, x3)), gepNor(x1, x2))))))
	}},
	"Id":    BoolNode{10019, "Id", 1, func(x0, x1, x2, x3 bool) bool { return (x0) }},
	"Zero":  BoolNode{10020, "Zero", 1, func(x0, x1, x2, x3 bool) bool { return (false) }},
	"One":   BoolNode{10021, "One", 1, func(x0, x1, x2, x3 bool) bool { return (true) }},
	"LT":    BoolNode{10022, "LT", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, gepNor(x1, x1)) }},
	"GT":    BoolNode{10023, "GT", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x0), x1) }},
	"LOE":   BoolNode{10024, "LOE", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x1), x1), (false)) }},
	"GOE":   BoolNode{10025, "GOE", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(x0, gepNor(x1, x1))) }},
	"NotA":  BoolNode{10026, "NotA", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, x0) }},
	"NotB":  BoolNode{10027, "NotB", 2, func(x0, x1, x2, x3 bool) bool { return gepNor(x1, x1) }},
	"IdA":   BoolNode{10028, "IdA", 2, func(x0, x1, x2, x3 bool) bool { return (x0) }},
	"IdB":   BoolNode{10029, "IdB", 2, func(x0, x1, x2, x3 bool) bool { return (x1) }},
	"Zero2": BoolNode{10030, "Zero2", 2, func(x0, x1, x2, x3 bool) bool { return (false) }},
	"One2":  BoolNode{10031, "One2", 2, func(x0, x1, x2, x3 bool) bool { return (true) }},
	"LT3":   BoolNode{10032, "LT3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x1)), gepNor(x2, x2)) }},
	"GT3":   BoolNode{10033, "GT3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x2), gepNor(x1, x2))) }},
	"LOE3":  BoolNode{10034, "LOE3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x2), gepNor(gepNor(x1, x2), x2)) }},
	"GOE3": BoolNode{10035, "GOE3", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, gepNor(gepNor(x0, x2), gepNor(x1, x1)))), (false))
	}},
	"Mux": BoolNode{10036, "Mux", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), x2)) }},
	"If":  BoolNode{10037, "If", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x1), x1), gepNor(x0, x2)) }},
	"Maj": BoolNode{10038, "Maj", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(gepNor(x0, x1), gepNor(x0, x2))), gepNor(x1, x2))
	}},
	"Min": BoolNode{10039, "Min", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x2)), gepNor(gepNor(x0, x2), gepNor(x1, x1)))
	}},
	"2Off": BoolNode{10040, "2Off", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), x2), gepNor(x0, gepNor(x0, x2))), gepNor(gepNor(x0, x2), gepNor(x1, x1)))
	}},
	"2On": BoolNode{10041, "2On", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x1, x1))), gepNor(x2, x2)), gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), x2))
	}},
	"LM3A1": BoolNode{10042, "LM3A1", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x2), gepNor(x1, gepNor(x2, x2))) }},
	"LM3A2": BoolNode{10043, "LM3A2", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x2), gepNor(gepNor(x1, x1), gepNor(x2, x2))) }},
	"LM3A3": BoolNode{10044, "LM3A3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x0), x2), gepNor(x1, gepNor(x1, x2))) }},
	"LM3A4": BoolNode{10045, "LM3A4", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x1)), gepNor(gepNor(x0, x2), x2))
	}},
	"LM3B1": BoolNode{10046, "LM3B1", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x2, x2)), gepNor(x1, x2)) }},
	"LM3B2": BoolNode{10047, "LM3B2", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x2)), gepNor(gepNor(x1, x2), x2)) }},
	"LM3B3": BoolNode{10048, "LM3B3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x0), gepNor(x2, x2)), gepNor(x1, x2)) }},
	"LM3B4": BoolNode{10049, "LM3B4", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x2, x2)), gepNor(gepNor(x1, x1), x2))
	}},
	"LM3C1": BoolNode{10050, "LM3C1", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x1), gepNor(gepNor(x1, x2), x2)) }},
	"LM3C2": BoolNode{10051, "LM3C2", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x1), gepNor(gepNor(x1, x1), gepNor(x2, x2))) }},
	"LM3C3": BoolNode{10052, "LM3C3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x1), x1), gepNor(gepNor(x1, x2), x2)) }},
	"LM3C4": BoolNode{10053, "LM3C4", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x2)), gepNor(gepNor(x0, x0), x1))
	}},
	"LM3D1": BoolNode{10054, "LM3D1", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x1)), gepNor(x1, x2)) }},
	"LM3D2": BoolNode{10055, "LM3D2", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x2)), gepNor(x1, gepNor(x1, x2))) }},
	"LM3D3": BoolNode{10056, "LM3D3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(x1, x2)) }},
	"LM3D4": BoolNode{10057, "LM3D4", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x2)), gepNor(x1, gepNor(x2, x2)))
	}},
	"LM3E1": BoolNode{10058, "LM3E1", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x0), gepNor(x1, x1))) }},
	"LM3E2": BoolNode{10059, "LM3E2", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(gepNor(x0, x2), x1)) }},
	"LM3E3": BoolNode{10060, "LM3E3", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(gepNor(x0, x2), gepNor(x1, x1)))
	}},
	"LM3F1": BoolNode{10061, "LM3F1", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(gepNor(x0, x1), x2)) }},
	"LM3F2": BoolNode{10062, "LM3F2", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x2, x2))) }},
	"LM3F3": BoolNode{10063, "LM3F3", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(gepNor(x0, x1), gepNor(x2, x2)))
	}},
	"LM3G1": BoolNode{10064, "LM3G1", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), x1), gepNor(x0, x2)), gepNor(x1, x2))
	}},
	"LM3G2": BoolNode{10065, "LM3G2", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x2, x2))), gepNor(gepNor(x0, x2), x1))
	}},
	"LM3G3": BoolNode{10066, "LM3G3", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x1, x2), x2)), gepNor(x0, x2))
	}},
	"LM3G4": BoolNode{10067, "LM3G4", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(x1, x2)))
	}},
	"LM3H1": BoolNode{10068, "LM3H1", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(x1, gepNor(x2, x2))) }},
	"LM3H2": BoolNode{10069, "LM3H2", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x0), x1), gepNor(gepNor(x0, x0), x2)) }},
	"LM3H3": BoolNode{10070, "LM3H3", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x1), x2) }},
	"LM3H4": BoolNode{10071, "LM3H4", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, gepNor(x1, x2)) }},
	"LT3A": BoolNode{10072, "LT3A", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(gepNor(x1, x1), gepNor(x2, x2)))
	}},
	"GT3A": BoolNode{10073, "GT3A", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x2), gepNor(x1, x2))), gepNor(gepNor(x1, x1), gepNor(x2, x2)))
	}},
	"LOE3A": BoolNode{10074, "LOE3A", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x0), gepNor(gepNor(x1, x1), x2)) }},
	"GOE3A": BoolNode{10075, "GOE3A", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x1, x2)), x2), gepNor(x0, gepNor(x2, x2)))
	}},
	"ET3A": BoolNode{10076, "ET3A", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x0), gepNor(gepNor(x1, x2), x2))))
	}},
	"NET3A": BoolNode{10077, "NET3A", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), gepNor(x1, x2))), gepNor(gepNor(x1, x1), gepNor(x2, x2)))
	}},
	"LT3B": BoolNode{10078, "LT3B", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x2), x2), gepNor(x1, x1)) }},
	"GT3B": BoolNode{10079, "GT3B", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), x1), x2), gepNor(x1, gepNor(x1, x2)))
	}},
	"LOE3B": BoolNode{10080, "LOE3B", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), x1), gepNor(gepNor(x0, x0), gepNor(x2, x2)))
	}},
	"GOE3B": BoolNode{10081, "GOE3B", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x1), gepNor(x0, x2))), gepNor(gepNor(x0, x0), gepNor(x2, x2)))
	}},
	"ET3B": BoolNode{10082, "ET3B", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), gepNor(x1, x1))), gepNor(gepNor(x0, x0), gepNor(x2, x2)))
	}},
	"NET3B": BoolNode{10083, "NET3B", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(gepNor(gepNor(x0, x0), x2), gepNor(x1, x1)), gepNor(x1, x2)))
	}},
	"LT3C": BoolNode{10084, "LT3C", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x2)), gepNor(x1, x2))), gepNor(x1, x2))
	}},
	"GT3C": BoolNode{10085, "GT3C", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x1), x1), gepNor(x1, x2)) }},
	"LOE3C": BoolNode{10086, "LOE3C", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x0), x1), gepNor(x0, x2)), gepNor(gepNor(x0, x2), x2))
	}},
	"GOE3C": BoolNode{10087, "GOE3C", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, gepNor(gepNor(x0, x2), x1)) }},
	"ET3C": BoolNode{10088, "ET3C", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x2)), gepNor(x0, x2)), gepNor(gepNor(x0, x0), gepNor(x1, x1)))
	}},
	"NET3C": BoolNode{10089, "NET3C", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(x1, x1))), gepNor(x1, x2))
	}},
	"T004": BoolNode{10090, "T004", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x2)), gepNor(x1, x2)) }},
	"T008": BoolNode{10091, "T008", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, gepNor(x0, gepNor(gepNor(x0, x1), gepNor(x0, x2)))) }},
	"T009": BoolNode{10092, "T009", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(x0, gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x2)), gepNor(x1, x2)))
	}},
	"T032": BoolNode{10093, "T032", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x2, x2)), x1), x1) }},
	"T033": BoolNode{10094, "T033", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x1, x2)), gepNor(x0, x2)), x1)
	}},
	"T041": BoolNode{10095, "T041", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x2)), gepNor(x1, x2)), gepNor(gepNor(x0, x0), gepNor(x1, gepNor(x2, x2))))
	}},
	"T055": BoolNode{10096, "T055", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x1)), (false)) }},
	"T057": BoolNode{10097, "T057", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x2)), gepNor(x1, x2)), gepNor(gepNor(x0, x0), gepNor(x1, x1)))
	}},
	"T064": BoolNode{10098, "T064", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x1)), x2), x2) }},
	"T065": BoolNode{10099, "T065", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), gepNor(x1, x2))), x2)
	}},
	"T069": BoolNode{10100, "T069", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x1), x1), x2) }},
	"T073": BoolNode{10101, "T073", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x1)), gepNor(x1, x2)), gepNor(gepNor(x0, x0), gepNor(gepNor(x1, x2), x2)))
	}},
	"T081": BoolNode{10102, "T081", 3, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x1)), x2) }},
	"T089": BoolNode{10103, "T089", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x1)), gepNor(x1, x2)), gepNor(gepNor(x0, x0), gepNor(x2, x2)))
	}},
	"T093": BoolNode{10104, "T093", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x2, x2)), gepNor(x1, gepNor(x1, x2)))
	}},
	"T096": BoolNode{10105, "T096", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x0), gepNor(gepNor(x1, gepNor(x1, x2)), gepNor(gepNor(x1, x2), x2)))
	}},
	"T101": BoolNode{10106, "T101", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x1, x2)), gepNor(x0, x2)), gepNor(gepNor(x1, x1), gepNor(x2, x2)))
	}},
	"T109": BoolNode{10107, "T109", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, gepNor(x2, x2)), x1), x1), gepNor(gepNor(x0, x0), gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x2)), gepNor(x1, x2))))
	}},
	"T111": BoolNode{10108, "T111", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x2)), gepNor(x1, x2))), (false))
	}},
	"T121": BoolNode{10109, "T121", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x2)), gepNor(x1, x2)), gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x0), gepNor(x1, x1))), gepNor(x2, x2)))
	}},
	"T123": BoolNode{10110, "T123", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x0), gepNor(x2, x2))), gepNor(x1, x1)), (false))
	}},
	"T125": BoolNode{10111, "T125", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x1, x1))), gepNor(x2, x2)), (false))
	}},
	"T154": BoolNode{10112, "T154", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), x1), gepNor(x2, x2)), gepNor(x1, x2)), gepNor(x0, x2))
	}},
	"T223": BoolNode{10113, "T223", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x2)), x1), x1), (false))
	}},
	"T239": BoolNode{10114, "T239", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x2), gepNor(x1, x2))), (false))
	}},
	"T249": BoolNode{10115, "T249", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x2)), gepNor(x1, x2)), (false))
	}},
	"T251": BoolNode{10116, "T251", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(x0, x2)), gepNor(x1, x1)), (false))
	}},
	"T253": BoolNode{10117, "T253", 3, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x2)), gepNor(x2, x2)), (false))
	}},
	"LT4": BoolNode{10118, "LT4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(x2, x2)), gepNor(x3, x3))
	}},
	"GT4": BoolNode{10119, "GT4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x2)), gepNor(gepNor(gepNor(x0, x3), x3), x2))
	}},
	"LOE4": BoolNode{10120, "LOE4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x3), gepNor(x1, x3))), gepNor(gepNor(x2, x2), x3))
	}},
	"GOE4": BoolNode{10121, "GOE4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(gepNor(x1, x1), gepNor(x2, x2))), gepNor(gepNor(x0, x2), gepNor(x0, gepNor(x3, x3)))), (false))
	}},
	"Tie": BoolNode{10122, "Tie", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x3)), gepNor(gepNor(x0, x3), gepNor(x1, x2))), gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(gepNor(x2, x2), gepNor(x3, x3))), gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x3)), gepNor(x0, x3))))
	}},
	"Ntie": BoolNode{10123, "Ntie", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), x2), gepNor(x0, gepNor(x0, x2))), gepNor(gepNor(x1, gepNor(x3, x3)), gepNor(gepNor(x1, x3), x3))), gepNor(gepNor(gepNor(gepNor(x0, x0), x3), gepNor(gepNor(x1, x2), x2)), gepNor(gepNor(x0, gepNor(x0, x3)), gepNor(x1, gepNor(x1, x2)))))
	}},
	"3Off": BoolNode{10124, "3Off", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), x2), x3), gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x1, x3)), x2), gepNor(x0, x3)))
	}},
	"3On": BoolNode{10125, "3On", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(x0, x1)), gepNor(gepNor(gepNor(x2, x2), gepNor(x3, x3)), gepNor(x2, x3))), gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(gepNor(x2, x2), gepNor(x3, x3))))
	}},
	"LM4A1": BoolNode{10126, "LM4A1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x3), gepNor(gepNor(gepNor(x1, x2), gepNor(x1, x3)), gepNor(x1, x3)))
	}},
	"LM4A2": BoolNode{10127, "LM4A2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x3), gepNor(gepNor(gepNor(gepNor(x1, x2), gepNor(x3, x3)), x2), x2))
	}},
	"LM4A3": BoolNode{10128, "LM4A3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x3), gepNor(x1, gepNor(x1, gepNor(gepNor(x1, x3), gepNor(x2, x2)))))
	}},
	"LM4A4": BoolNode{10129, "LM4A4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x3), gepNor(gepNor(gepNor(x0, x3), gepNor(gepNor(x1, x1), gepNor(x2, x2))), gepNor(gepNor(x1, x3), x3)))
	}},
	"LM4A5": BoolNode{10130, "LM4A5", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x3), gepNor(gepNor(x1, x2), gepNor(x1, x3))), gepNor(gepNor(x0, x0), x3))
	}},
	"LM4A6": BoolNode{10131, "LM4A6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x3), gepNor(x1, x1)), x2), x2), gepNor(gepNor(x0, x3), x3))
	}},
	"LM4A7": BoolNode{10132, "LM4A7", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x3), x3), gepNor(x1, gepNor(x1, gepNor(gepNor(x1, x2), gepNor(x1, x3)))))
	}},
	"LM4A8": BoolNode{10133, "LM4A8", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x3), gepNor(x2, x2))), gepNor(x1, x1)), gepNor(gepNor(x0, x3), x3))
	}},
	"LM4B1": BoolNode{10134, "LM4B1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x3, x3)), gepNor(gepNor(gepNor(x1, x3), gepNor(x2, x3)), x3))
	}},
	"LM4B2": BoolNode{10135, "LM4B2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x1, x3)), gepNor(gepNor(x1, gepNor(x2, x3)), x3))
	}},
	"LM4B3": BoolNode{10136, "LM4B3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x3)), gepNor(gepNor(gepNor(x1, x3), x2), x3))
	}},
	"LM4B4": BoolNode{10137, "LM4B4", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x3)), gepNor(gepNor(x1, x2), x3)) }},
	"LM4B5": BoolNode{10138, "LM4B5", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x3, x3)), gepNor(gepNor(gepNor(x1, x3), gepNor(x2, x3)), gepNor(x2, gepNor(x3, x3))))
	}},
	"LM4B6": BoolNode{10139, "LM4B6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x3)), gepNor(gepNor(x1, gepNor(x1, x2)), x3))
	}},
	"LM4B7": BoolNode{10140, "LM4B7", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x2, x3)), gepNor(gepNor(gepNor(x1, x3), x2), x3))
	}},
	"LM4B8": BoolNode{10141, "LM4B8", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x3)), gepNor(gepNor(x1, x2), x3))
	}},
	"LM4C1": BoolNode{10142, "LM4C1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x3)), gepNor(x1, x2)))
	}},
	"LM4C2": BoolNode{10143, "LM4C2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(gepNor(gepNor(x1, x3), gepNor(x2, x3)), x3), x3))
	}},
	"LM4C3": BoolNode{10144, "LM4C3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x3)), x2), x2))
	}},
	"LM4C4": BoolNode{10145, "LM4C4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x2)), gepNor(x2, x3)), gepNor(x3, x3)))
	}},
	"LM4C5": BoolNode{10146, "LM4C5", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), x1), gepNor(gepNor(x1, x2), gepNor(gepNor(x2, x3), gepNor(x2, x3))))
	}},
	"LM4C6": BoolNode{10147, "LM4C6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), x1), gepNor(gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x2)), x3), x3))
	}},
	"LM4C7": BoolNode{10148, "LM4C7", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x3)), x2), x2), gepNor(gepNor(x0, x1), x1))
	}},
	"LM4C8": BoolNode{10149, "LM4C8", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), x1), gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x2)), gepNor(x2, x3)), gepNor(x3, x3)))
	}},
	"LM4D1": BoolNode{10150, "LM4D1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(x1, gepNor(gepNor(x1, x2), gepNor(x3, x3))))
	}},
	"LM4D2": BoolNode{10151, "LM4D2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x1, x2)), gepNor(x1, gepNor(x2, gepNor(x3, x3))))
	}},
	"LM4D3": BoolNode{10152, "LM4D3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x1, x3)), gepNor(x1, gepNor(gepNor(x2, x3), x3)))
	}},
	"LM4D4": BoolNode{10153, "LM4D4", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x3)), gepNor(x1, gepNor(x2, x3))) }},
	"LM4D5": BoolNode{10154, "LM4D5", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(x1, gepNor(gepNor(x1, x3), gepNor(x2, x2))))
	}},
	"LM4D6": BoolNode{10155, "LM4D6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x1, x1)), gepNor(x1, x2)), gepNor(x1, x3))
	}},
	"LM4D7": BoolNode{10156, "LM4D7", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(x1, x3)), gepNor(x1, x2))
	}},
	"LM4D8": BoolNode{10157, "LM4D8", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x2)), gepNor(x1, gepNor(x2, x3)))
	}},
	"LM4E1": BoolNode{10158, "LM4E1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x2), gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x3)), gepNor(x1, x2)))
	}},
	"LM4E2": BoolNode{10159, "LM4E2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x2), gepNor(gepNor(gepNor(gepNor(x1, x3), gepNor(x2, x3)), x3), x3))
	}},
	"LM4E3": BoolNode{10160, "LM4E3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x2), gepNor(x1, gepNor(x1, gepNor(gepNor(x1, x2), gepNor(x3, x3)))))
	}},
	"LM4E4": BoolNode{10161, "LM4E4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x2, x2), gepNor(x3, x3))), gepNor(x1, x1)), gepNor(x0, x2))
	}},
	"LM4E5": BoolNode{10162, "LM4E5", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x2), x2), gepNor(gepNor(gepNor(gepNor(x0, x2), x1), x3), x3))
	}},
	"LM4E6": BoolNode{10163, "LM4E6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x3)), x3), x3), gepNor(gepNor(x0, x0), x2))
	}},
	"LM4E7": BoolNode{10164, "LM4E7", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), x2), gepNor(x1, gepNor(x1, gepNor(gepNor(x2, x2), gepNor(x3, x3)))))
	}},
	"LM4E8": BoolNode{10165, "LM4E8", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), x2), gepNor(gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x2)), gepNor(x2, x3)), gepNor(x3, x3)))
	}},
	"LM4F1": BoolNode{10166, "LM4F1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(gepNor(gepNor(x1, x2), gepNor(x2, x3)), x2))
	}},
	"LM4F2": BoolNode{10167, "LM4F2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x1, x2)), gepNor(gepNor(x1, gepNor(x2, x3)), x2))
	}},
	"LM4F3": BoolNode{10168, "LM4F3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(gepNor(gepNor(x1, x1), x3), x2))
	}},
	"LM4F4": BoolNode{10169, "LM4F4", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x2)), gepNor(gepNor(x1, x3), x2)) }},
	"LM4F5": BoolNode{10170, "LM4F5", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x2, x2)), gepNor(gepNor(x1, x1), x2)), gepNor(x2, x3))
	}},
	"LM4F6": BoolNode{10171, "LM4F6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(x1, x2)), gepNor(x2, x3))
	}},
	"LM4F7": BoolNode{10172, "LM4F7", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x2, x2)), gepNor(x2, x3)), gepNor(x1, x2))
	}},
	"LM4F8": BoolNode{10173, "LM4F8", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(x2, x3)), gepNor(gepNor(x1, x3), x2))
	}},
	"LM4G1": BoolNode{10174, "LM4G1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x2), x2), x1), x1), gepNor(x0, x3))
	}},
	"LM4G2": BoolNode{10175, "LM4G2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(x1, x2)), x2), x2), gepNor(x0, x3))
	}},
	"LM4G3": BoolNode{10176, "LM4G3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x2, x2)), x1), x1), gepNor(x0, x3))
	}},
	"LM4G4": BoolNode{10177, "LM4G4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(x2, x3)), gepNor(x2, x2)), gepNor(x0, x3))
	}},
	"LM4G5": BoolNode{10178, "LM4G5", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x3, x3)), gepNor(gepNor(gepNor(gepNor(x0, x3), x2), x1), x1))
	}},
	"LM4G6": BoolNode{10179, "LM4G6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x3)), gepNor(gepNor(gepNor(gepNor(x0, x3), gepNor(x1, x2)), x2), x2))
	}},
	"LM4G7": BoolNode{10180, "LM4G7", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x3)), gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x2)), x1), x1))
	}},
	"LM4G8": BoolNode{10181, "LM4G8", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x3), gepNor(x1, x1))), gepNor(x2, x2)), gepNor(x0, gepNor(x3, x3)))
	}},
	"LM4H1": BoolNode{10182, "LM4H1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(gepNor(x0, x2), gepNor(x1, x1))), gepNor(gepNor(x0, x3), x3))
	}},
	"LM4H2": BoolNode{10183, "LM4H2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(gepNor(x0, x2), x1)), gepNor(gepNor(x0, x1), x3))
	}},
	"LM4H3": BoolNode{10184, "LM4H3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(gepNor(x0, x1), x2)), gepNor(gepNor(x0, x0), x3))
	}},
	"LM4H4": BoolNode{10185, "LM4H4", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x2)), gepNor(gepNor(x0, x1), x3)) }},
	"LM4H5": BoolNode{10186, "LM4H5", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(gepNor(x0, x3), x3)))
	}},
	"LM4H6": BoolNode{10187, "LM4H6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(gepNor(x1, x2), x1)), gepNor(gepNor(x0, x0), gepNor(x3, x3)))
	}},
	"LM4H7": BoolNode{10188, "LM4H7", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x3), x3)))
	}},
	"LM4H8": BoolNode{10189, "LM4H8", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x1, x2)), gepNor(gepNor(x0, x2), gepNor(x3, x3)))
	}},
	"LM4I1": BoolNode{10190, "LM4I1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(x2, x2)), gepNor(x3, x3)), gepNor(x2, gepNor(x3, x3)))
	}},
	"LM4I2": BoolNode{10191, "LM4I2", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), x1), gepNor(gepNor(x0, x0), gepNor(gepNor(x0, x2), gepNor(gepNor(x2, x2), gepNor(x3, x3)))))
	}},
	"LM4I3": BoolNode{10192, "LM4I3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(gepNor(x1, x1), gepNor(x2, x2))), gepNor(x3, x3)), (false))
	}},
	"LM4I4": BoolNode{10193, "LM4I4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(gepNor(x0, x1), gepNor(gepNor(x1, x1), gepNor(x2, x2)))), gepNor(gepNor(x0, x3), x3))
	}},
	"LM4I5": BoolNode{10194, "LM4I5", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(gepNor(x0, x1), x2), x3) }},
	"LM4I6": BoolNode{10195, "LM4I6", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, gepNor(x1, gepNor(x2, x3))) }},
	"LM4I7": BoolNode{10196, "LM4I7", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, gepNor(x1, x2)), x3) }},
	"LM4I8": BoolNode{10197, "LM4I8", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(x0, gepNor(gepNor(x1, x2), x3)) }},
	"LT4A": BoolNode{10198, "LT4A", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(gepNor(gepNor(x0, x1), x2), gepNor(x1, x3))), gepNor(gepNor(x0, x3), x3))
	}},
	"GT4A": BoolNode{10199, "GT4A", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), x1), x2), x2), gepNor(gepNor(gepNor(x0, x3), x1), x3))
	}},
	"LOE4A": BoolNode{10200, "LOE4A", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(gepNor(gepNor(x1, x1), x2), gepNor(x1, x3))), gepNor(x0, x2))
	}},
	"GOE4A": BoolNode{10201, "GOE4A", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x1, x1)), x2), gepNor(gepNor(gepNor(x0, gepNor(x0, x1)), x3), x3))
	}},
	"ET4A": BoolNode{10202, "ET4A", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), gepNor(x1, x1))), x2), gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x3), gepNor(x1, x3))))
	}},
	"NET4A": BoolNode{10203, "NET4A", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x1, x3))), x3), gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), gepNor(x1, x2))))
	}},
	"LT4B": BoolNode{10204, "LT4B", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x3)), gepNor(x1, x3)), gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(x3, x3)))
	}},
	"GT4B": BoolNode{10205, "GT4B", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x3)), gepNor(x0, x3)), gepNor(gepNor(x1, gepNor(x1, x2)), gepNor(x3, x3)))
	}},
	"LOE4B": BoolNode{10206, "LOE4B", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), x1), x3), gepNor(gepNor(gepNor(gepNor(x0, x0), x1), gepNor(x2, x2)), gepNor(x1, x3)))
	}},
	"GOE4B": BoolNode{10207, "GOE4B", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(x2, x2)), gepNor(x0, x3)), gepNor(gepNor(x0, gepNor(x1, x3)), x3))
	}},
	"ET4B": BoolNode{10208, "ET4B", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x0, gepNor(x0, x1))), gepNor(gepNor(x0, gepNor(x0, x3)), gepNor(x1, gepNor(x1, x3)))), gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x3), gepNor(x1, x3))), gepNor(gepNor(x2, x2), gepNor(x3, x3))))
	}},
	"NET4B": BoolNode{10209, "NET4B", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x3)), gepNor(x1, x3)), gepNor(gepNor(gepNor(gepNor(x0, x1), x1), x3), gepNor(gepNor(gepNor(x0, x0), gepNor(x1, gepNor(x1, x3))), gepNor(x2, x2))))
	}},
	"LT4C": BoolNode{10210, "LT4C", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(gepNor(x1, x2), x3)), gepNor(gepNor(x0, x2), x3)), gepNor(x1, gepNor(x1, x3)))
	}},
	"GT4C": BoolNode{10211, "GT4C", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x3), x1), x3), gepNor(x1, gepNor(x2, x3))), gepNor(x0, gepNor(x0, x3)))
	}},
	"LOE4C": BoolNode{10212, "LOE4C", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x3, x3)), x1), x1), gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x1, x3), x2)), x3))
	}},
	"GOE4C": BoolNode{10213, "GOE4C", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(x2, x3)), gepNor(x0, x3)), gepNor(gepNor(gepNor(x1, x2), x3), x3))
	}},
	"ET4C": BoolNode{10214, "ET4C", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x3)), x1), x1), gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x3), gepNor(x1, x3))), gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(x2, x3))))
	}},
	"NET4C": BoolNode{10215, "NET4C", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x1, x1))), gepNor(x3, x3)), gepNor(gepNor(gepNor(gepNor(x0, x1), x2), x3), gepNor(gepNor(x0, x3), gepNor(x1, x3))))
	}},
	"LT4D": BoolNode{10216, "LT4D", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x1, x1)), gepNor(gepNor(x2, x2), gepNor(x3, x3)))
	}},
	"GT4D": BoolNode{10217, "GT4D", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(gepNor(x1, x1), gepNor(x2, x2))), gepNor(x3, x3)), gepNor(gepNor(gepNor(gepNor(x0, x3), x1), x3), x3))
	}},
	"LOE4D": BoolNode{10218, "LOE4D", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x0), gepNor(gepNor(gepNor(x1, gepNor(x1, x2)), gepNor(x3, x3)), gepNor(x1, x3)))
	}},
	"GOE4D": BoolNode{10219, "GOE4D", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(x3, x3)), gepNor(x0, x3)))
	}},
	"ET4D": BoolNode{10220, "ET4D", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), gepNor(x3, x3)))), gepNor(gepNor(gepNor(gepNor(x0, x3), x3), gepNor(x1, gepNor(x1, x2))), gepNor(x1, x3)))
	}},
	"NET4D": BoolNode{10221, "NET4D", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(gepNor(gepNor(gepNor(x0, x1), x1), x3), gepNor(gepNor(gepNor(gepNor(x0, x1), x1), gepNor(x2, x2)), gepNor(x3, x3))))
	}},
	"LT4E": BoolNode{10222, "LT4E", 4, func(x0, x1, x2, x3 bool) bool { return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x2), x2)) }},
	"GT4E": BoolNode{10223, "GT4E", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x0), gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x2)), gepNor(x1, x3)))
	}},
	"LOE4E": BoolNode{10224, "LOE4E", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(gepNor(gepNor(x1, x1), gepNor(x3, x3)), gepNor(x1, x2))), (false))
	}},
	"GOE4E": BoolNode{10225, "GOE4E", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x1, x1)), gepNor(gepNor(x0, x1), gepNor(x3, x3)))
	}},
	"ET4E": BoolNode{10226, "ET4E", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x3), gepNor(x1, x2))), gepNor(gepNor(x1, x1), gepNor(x3, x3)))
	}},
	"NET4E": BoolNode{10227, "NET4E", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(gepNor(gepNor(x0, x0), x2), gepNor(x1, x1)), gepNor(x1, x3)))
	}},
	"Q0002": BoolNode{10228, "Q0002", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(x0, gepNor(x0, gepNor(gepNor(x0, gepNor(x1, x2)), gepNor(x2, x3))))
	}},
	"Q001C": BoolNode{10229, "Q001C", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x3)), gepNor(gepNor(x0, gepNor(x1, x2)), x2))
	}},
	"Q0048": BoolNode{10230, "Q0048", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x3)), gepNor(x0, gepNor(gepNor(x1, x3), gepNor(x2, x2))))
	}},
	"Q0800": BoolNode{10231, "Q0800", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x2)), x1), gepNor(x1, x3)), x1), x1)
	}},
	"Q3378": BoolNode{10232, "Q3378", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), x2), gepNor(x0, gepNor(gepNor(x0, x1), gepNor(x3, x3)))), gepNor(gepNor(gepNor(x0, x2), gepNor(x0, x3)), x1))
	}},
	"Q3475": BoolNode{10233, "Q3475", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), x3), gepNor(x2, x2)), gepNor(gepNor(gepNor(gepNor(x0, x3), x2), x1), x1))
	}},
	"Q3CB0": BoolNode{10234, "Q3CB0", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x1), gepNor(gepNor(x0, gepNor(x3, x3)), gepNor(gepNor(gepNor(x1, x1), gepNor(x2, x2)), gepNor(x1, x2))))
	}},
	"Q3DEF": BoolNode{10235, "Q3DEF", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), x1), gepNor(gepNor(gepNor(x0, x2), x2), gepNor(x0, gepNor(x2, x3)))), gepNor(x1, gepNor(x2, gepNor(x2, x3)))), (false))
	}},
	"Q3DFF": BoolNode{10236, "Q3DFF", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), x2), x1), gepNor(gepNor(gepNor(x0, x0), gepNor(x2, x2)), gepNor(x1, gepNor(x2, x3)))), (false))
	}},
	"Q4200": BoolNode{10237, "Q4200", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x3)), gepNor(gepNor(x0, x2), gepNor(x3, x3))), gepNor(gepNor(x1, x2), gepNor(gepNor(x2, x3), x3)))
	}},
	"Q4C11": BoolNode{10238, "Q4C11", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), x1), gepNor(x3, x3)), gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x3), gepNor(x2, x2))))
	}},
	"Q5100": BoolNode{10239, "Q5100", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, x3), gepNor(gepNor(gepNor(x1, x1), x3), gepNor(x2, x3)))
	}},
	"Q5EEF": BoolNode{10240, "Q5EEF", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x0), gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(gepNor(x1, x1), gepNor(x3, x3)))), gepNor(gepNor(gepNor(gepNor(x0, x1), x3), x2), gepNor(gepNor(x0, x2), gepNor(x1, x2))))
	}},
	"Q5EFF": BoolNode{10241, "Q5EFF", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x1)), gepNor(gepNor(x0, x2), x3)), gepNor(gepNor(x1, x2), x3)), (false))
	}},
	"Q6A6D": BoolNode{10242, "Q6A6D", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), x2), gepNor(gepNor(x0, x0), x1)), x3), gepNor(gepNor(gepNor(x0, x1), x2), gepNor(gepNor(gepNor(x1, x1), gepNor(x3, x3)), gepNor(gepNor(x1, x3), x2))))
	}},
	"Q6F75": BoolNode{10243, "Q6F75", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x3), x3), gepNor(gepNor(x1, x1), gepNor(x2, x2))), gepNor(gepNor(gepNor(x1, x1), x2), x3)), gepNor(x0, gepNor(x0, gepNor(x1, gepNor(x3, x3)))))
	}},
	"Q74C4": BoolNode{10244, "Q74C4", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, gepNor(x0, x1)), gepNor(gepNor(x1, x2), gepNor(gepNor(x2, x2), gepNor(x3, x3)))), gepNor(x0, x2))
	}},
	"Q7DA3": BoolNode{10245, "Q7DA3", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(x0, gepNor(x0, gepNor(x1, gepNor(x2, x2)))), gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x3, x3)), gepNor(gepNor(gepNor(x0, x3), x3), x3)), gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x0, x0), gepNor(x1, x2))), gepNor(gepNor(x1, x1), gepNor(x2, x2)))))
	}},
	"Q8304": BoolNode{10246, "Q8304", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, gepNor(x0, x3)), gepNor(gepNor(x0, x3), x1)), gepNor(x1, x2)), gepNor(gepNor(gepNor(x0, x2), x1), gepNor(gepNor(x2, x2), gepNor(x3, x3))))
	}},
	"Q8430": BoolNode{10247, "Q8430", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x3, x3)), gepNor(x1, gepNor(x2, x3))), gepNor(x2, x2))), gepNor(gepNor(gepNor(x0, x3), x3), x1))
	}},
	"Q8543": BoolNode{10248, "Q8543", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, gepNor(x0, x2)), gepNor(gepNor(x0, x2), x1)), gepNor(gepNor(x0, x1), x3)), gepNor(gepNor(gepNor(gepNor(x0, x0), x3), gepNor(gepNor(x1, x2), x2)), gepNor(x1, x1)))
	}},
	"Q9D80": BoolNode{10249, "Q9D80", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(gepNor(x1, x1), x3)), gepNor(gepNor(x0, x3), x2)), gepNor(x2, gepNor(x2, x3)))
	}},
	"QA092": BoolNode{10250, "QA092", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x2)), gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x3)), gepNor(x1, x3))), gepNor(gepNor(x0, x2), gepNor(gepNor(x1, x1), gepNor(x3, x3))))
	}},
	"QB36A": BoolNode{10251, "QB36A", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), x0), gepNor(x2, x3)), gepNor(x0, x2)), x3), gepNor(gepNor(x0, gepNor(gepNor(x0, x1), gepNor(x0, x2))), gepNor(gepNor(x0, gepNor(x3, x3)), gepNor(x1, gepNor(x1, x2)))))
	}},
	"QCBCF": BoolNode{10252, "QCBCF", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x1, x2)), x1), gepNor(gepNor(x1, x3), gepNor(x2, x3))), gepNor(gepNor(x1, x1), x2))
	}},
	"QEEB1": BoolNode{10253, "QEEB1", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x3), gepNor(x2, x3))), gepNor(gepNor(gepNor(x0, x1), gepNor(x2, x3)), gepNor(x0, x3)))
	}},
	"QEFFF": BoolNode{10254, "QEFFF", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x0), gepNor(x1, x3)), x3), gepNor(gepNor(x0, x3), gepNor(x2, x3))), (false))
	}},
	"QFF7B": BoolNode{10255, "QFF7B", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, gepNor(x0, x3))), gepNor(gepNor(gepNor(x0, x2), x3), gepNor(gepNor(x1, x1), gepNor(x2, x2)))), (false))
	}},
	"QFFF6": BoolNode{10256, "QFFF6", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(x0, x1), gepNor(x0, x1)), gepNor(gepNor(gepNor(x1, x2), gepNor(x3, x3)), gepNor(x2, x3))), (false))
	}},
	"QFFFB": BoolNode{10257, "QFFFB", 4, func(x0, x1, x2, x3 bool) bool {
		return gepNor(gepNor(gepNor(gepNor(gepNor(x0, x2), gepNor(gepNor(x0, x3), x1)), x1), x1), (false))
	}},
}
