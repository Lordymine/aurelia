// index.ts
import { createInterface } from "node:readline";

// node_modules/@anthropic-ai/claude-agent-sdk/sdk.mjs
import { join as HV } from "path";
import { fileURLToPath as lP } from "url";
import { setMaxListeners as _V } from "events";
import { spawn as ZU } from "child_process";
import { createInterface as CU } from "readline";
import { randomUUID as nq } from "crypto";
import { appendFile as oq, mkdir as rq } from "fs/promises";
import { join as g5 } from "path";
import { join as dq } from "path";
import { homedir as iq } from "os";
import { appendFile as UU, mkdir as LU, unlink as FU, symlink as NU } from "fs/promises";
import { dirname as n5, join as o5 } from "path";
import { cwd as aq } from "process";
import { realpathSync as h5 } from "fs";
import { randomUUID as n9 } from "crypto";
import * as u from "fs";
import { stat as XU, readdir as YU, readFile as d5, unlink as JU, rmdir as GU, rm as WU, mkdir as HU, rename as BU, open as zU } from "fs/promises";
import { execFile as TU } from "child_process";
import { promisify as xU } from "util";
var MV = Object.create;
var { getPrototypeOf: AV, defineProperty: fQ, getOwnPropertyNames: jV } = Object;
var RV = Object.prototype.hasOwnProperty;
function IV(Q) {
  return this[Q];
}
var EV;
var PV;
var G5 = (Q, $, X) => {
  var Y = Q != null && typeof Q === "object";
  if (Y) {
    var J = $ ? EV ??= /* @__PURE__ */ new WeakMap() : PV ??= /* @__PURE__ */ new WeakMap(), G = J.get(Q);
    if (G) return G;
  }
  X = Q != null ? MV(AV(Q)) : {};
  let W = $ || !Q || !Q.__esModule ? fQ(X, "default", { value: Q, enumerable: true }) : X;
  for (let H of jV(Q)) if (!RV.call(W, H)) fQ(W, H, { get: IV.bind(Q, H), enumerable: true });
  if (Y) J.set(Q, W);
  return W;
};
var P = (Q, $) => () => ($ || Q(($ = { exports: {} }).exports, $), $.exports);
var bV = (Q) => Q;
function ZV(Q, $) {
  this[Q] = bV.bind(null, $);
}
var uQ = (Q, $) => {
  for (var X in $) fQ(Q, X, { get: $[X], enumerable: true, configurable: true, set: ZV.bind($, X) });
};
var CV = Symbol.dispose || /* @__PURE__ */ Symbol.for("Symbol.dispose");
var SV = Symbol.asyncDispose || /* @__PURE__ */ Symbol.for("Symbol.asyncDispose");
var $0 = (Q, $, X) => {
  if ($ != null) {
    if (typeof $ !== "object" && typeof $ !== "function") throw TypeError('Object expected to be assigned to "using" declaration');
    var Y;
    if (X) Y = $[SV];
    if (Y === void 0) Y = $[CV];
    if (typeof Y !== "function") throw TypeError("Object not disposable");
    Q.push([X, Y, $]);
  } else if (X) Q.push([X]);
  return $;
};
var X0 = (Q, $, X) => {
  var Y = typeof SuppressedError === "function" ? SuppressedError : function(W, H, B, z) {
    return z = Error(B), z.name = "SuppressedError", z.error = W, z.suppressed = H, z;
  }, J = (W) => $ = X ? new Y(W, $, "An error was suppressed during disposal") : (X = true, W), G = (W) => {
    while (W = Q.pop()) try {
      var H = W[1] && W[1].call(W[2]);
      if (W[0]) return Promise.resolve(H).then(G, (B) => (J(B), G()));
    } catch (B) {
      J(B);
    }
    if (X) throw $;
  };
  return G();
};
var M9 = P((U3) => {
  Object.defineProperty(U3, "__esModule", { value: true });
  U3.regexpCode = U3.getEsmExportName = U3.getProperty = U3.safeStringify = U3.stringify = U3.strConcat = U3.addCodeArg = U3.str = U3._ = U3.nil = U3._Code = U3.Name = U3.IDENTIFIER = U3._CodeOrName = void 0;
  class QQ {
  }
  U3._CodeOrName = QQ;
  U3.IDENTIFIER = /^[a-z$_][a-z$_0-9]*$/i;
  class w4 extends QQ {
    constructor(Q) {
      super();
      if (!U3.IDENTIFIER.test(Q)) throw Error("CodeGen: name must be a valid identifier");
      this.str = Q;
    }
    toString() {
      return this.str;
    }
    emptyStr() {
      return false;
    }
    get names() {
      return { [this.str]: 1 };
    }
  }
  U3.Name = w4;
  class z1 extends QQ {
    constructor(Q) {
      super();
      this._items = typeof Q === "string" ? [Q] : Q;
    }
    toString() {
      return this.str;
    }
    emptyStr() {
      if (this._items.length > 1) return false;
      let Q = this._items[0];
      return Q === "" || Q === '""';
    }
    get str() {
      var Q;
      return (Q = this._str) !== null && Q !== void 0 ? Q : this._str = this._items.reduce(($, X) => `${$}${X}`, "");
    }
    get names() {
      var Q;
      return (Q = this._names) !== null && Q !== void 0 ? Q : this._names = this._items.reduce(($, X) => {
        if (X instanceof w4) $[X.str] = ($[X.str] || 0) + 1;
        return $;
      }, {});
    }
  }
  U3._Code = z1;
  U3.nil = new z1("");
  function V3(Q, ...$) {
    let X = [Q[0]], Y = 0;
    while (Y < $.length) pY(X, $[Y]), X.push(Q[++Y]);
    return new z1(X);
  }
  U3._ = V3;
  var cY = new z1("+");
  function q3(Q, ...$) {
    let X = [w9(Q[0])], Y = 0;
    while (Y < $.length) X.push(cY), pY(X, $[Y]), X.push(cY, w9(Q[++Y]));
    return Rw(X), new z1(X);
  }
  U3.str = q3;
  function pY(Q, $) {
    if ($ instanceof z1) Q.push(...$._items);
    else if ($ instanceof w4) Q.push($);
    else Q.push(Pw($));
  }
  U3.addCodeArg = pY;
  function Rw(Q) {
    let $ = 1;
    while ($ < Q.length - 1) {
      if (Q[$] === cY) {
        let X = Iw(Q[$ - 1], Q[$ + 1]);
        if (X !== void 0) {
          Q.splice($ - 1, 3, X);
          continue;
        }
        Q[$++] = "+";
      }
      $++;
    }
  }
  function Iw(Q, $) {
    if ($ === '""') return Q;
    if (Q === '""') return $;
    if (typeof Q == "string") {
      if ($ instanceof w4 || Q[Q.length - 1] !== '"') return;
      if (typeof $ != "string") return `${Q.slice(0, -1)}${$}"`;
      if ($[0] === '"') return Q.slice(0, -1) + $.slice(1);
      return;
    }
    if (typeof $ == "string" && $[0] === '"' && !(Q instanceof w4)) return `"${Q}${$.slice(1)}`;
    return;
  }
  function Ew(Q, $) {
    return $.emptyStr() ? Q : Q.emptyStr() ? $ : q3`${Q}${$}`;
  }
  U3.strConcat = Ew;
  function Pw(Q) {
    return typeof Q == "number" || typeof Q == "boolean" || Q === null ? Q : w9(Array.isArray(Q) ? Q.join(",") : Q);
  }
  function bw(Q) {
    return new z1(w9(Q));
  }
  U3.stringify = bw;
  function w9(Q) {
    return JSON.stringify(Q).replace(/\u2028/g, "\\u2028").replace(/\u2029/g, "\\u2029");
  }
  U3.safeStringify = w9;
  function Zw(Q) {
    return typeof Q == "string" && U3.IDENTIFIER.test(Q) ? new z1(`.${Q}`) : V3`[${Q}]`;
  }
  U3.getProperty = Zw;
  function Cw(Q) {
    if (typeof Q == "string" && U3.IDENTIFIER.test(Q)) return new z1(`${Q}`);
    throw Error(`CodeGen: invalid export name: ${Q}, use explicit $id name mapping`);
  }
  U3.getEsmExportName = Cw;
  function Sw(Q) {
    return new z1(Q.toString());
  }
  U3.regexpCode = Sw;
});
var oY = P((O3) => {
  Object.defineProperty(O3, "__esModule", { value: true });
  O3.ValueScope = O3.ValueScopeName = O3.Scope = O3.varKinds = O3.UsedValueState = void 0;
  var g0 = M9();
  class F3 extends Error {
    constructor(Q) {
      super(`CodeGen: "code" for ${Q} not defined`);
      this.value = Q.value;
    }
  }
  var XQ;
  (function(Q) {
    Q[Q.Started = 0] = "Started", Q[Q.Completed = 1] = "Completed";
  })(XQ || (O3.UsedValueState = XQ = {}));
  O3.varKinds = { const: new g0.Name("const"), let: new g0.Name("let"), var: new g0.Name("var") };
  class iY {
    constructor({ prefixes: Q, parent: $ } = {}) {
      this._names = {}, this._prefixes = Q, this._parent = $;
    }
    toName(Q) {
      return Q instanceof g0.Name ? Q : this.name(Q);
    }
    name(Q) {
      return new g0.Name(this._newName(Q));
    }
    _newName(Q) {
      let $ = this._names[Q] || this._nameGroup(Q);
      return `${Q}${$.index++}`;
    }
    _nameGroup(Q) {
      var $, X;
      if (((X = ($ = this._parent) === null || $ === void 0 ? void 0 : $._prefixes) === null || X === void 0 ? void 0 : X.has(Q)) || this._prefixes && !this._prefixes.has(Q)) throw Error(`CodeGen: prefix "${Q}" is not allowed in this scope`);
      return this._names[Q] = { prefix: Q, index: 0 };
    }
  }
  O3.Scope = iY;
  class nY extends g0.Name {
    constructor(Q, $) {
      super($);
      this.prefix = Q;
    }
    setValue(Q, { property: $, itemIndex: X }) {
      this.value = Q, this.scopePath = g0._`.${new g0.Name($)}[${X}]`;
    }
  }
  O3.ValueScopeName = nY;
  var cw = g0._`\n`;
  class N3 extends iY {
    constructor(Q) {
      super(Q);
      this._values = {}, this._scope = Q.scope, this.opts = { ...Q, _n: Q.lines ? cw : g0.nil };
    }
    get() {
      return this._scope;
    }
    name(Q) {
      return new nY(Q, this._newName(Q));
    }
    value(Q, $) {
      var X;
      if ($.ref === void 0) throw Error("CodeGen: ref must be passed in value");
      let Y = this.toName(Q), { prefix: J } = Y, G = (X = $.key) !== null && X !== void 0 ? X : $.ref, W = this._values[J];
      if (W) {
        let z = W.get(G);
        if (z) return z;
      } else W = this._values[J] = /* @__PURE__ */ new Map();
      W.set(G, Y);
      let H = this._scope[J] || (this._scope[J] = []), B = H.length;
      return H[B] = $.ref, Y.setValue($, { property: J, itemIndex: B }), Y;
    }
    getValue(Q, $) {
      let X = this._values[Q];
      if (!X) return;
      return X.get($);
    }
    scopeRefs(Q, $ = this._values) {
      return this._reduceValues($, (X) => {
        if (X.scopePath === void 0) throw Error(`CodeGen: name "${X}" has no value`);
        return g0._`${Q}${X.scopePath}`;
      });
    }
    scopeCode(Q = this._values, $, X) {
      return this._reduceValues(Q, (Y) => {
        if (Y.value === void 0) throw Error(`CodeGen: name "${Y}" has no value`);
        return Y.value.code;
      }, $, X);
    }
    _reduceValues(Q, $, X = {}, Y) {
      let J = g0.nil;
      for (let G in Q) {
        let W = Q[G];
        if (!W) continue;
        let H = X[G] = X[G] || /* @__PURE__ */ new Map();
        W.forEach((B) => {
          if (H.has(B)) return;
          H.set(B, XQ.Started);
          let z = $(B);
          if (z) {
            let K = this.opts.es5 ? O3.varKinds.var : O3.varKinds.const;
            J = g0._`${J}${K} ${B} = ${z};${this.opts._n}`;
          } else if (z = Y === null || Y === void 0 ? void 0 : Y(B)) J = g0._`${J}${z}${this.opts._n}`;
          else throw new F3(B);
          H.set(B, XQ.Completed);
        });
      }
      return J;
    }
  }
  O3.ValueScope = N3;
});
var c = P((h0) => {
  Object.defineProperty(h0, "__esModule", { value: true });
  h0.or = h0.and = h0.not = h0.CodeGen = h0.operators = h0.varKinds = h0.ValueScopeName = h0.ValueScope = h0.Scope = h0.Name = h0.regexpCode = h0.stringify = h0.getProperty = h0.nil = h0.strConcat = h0.str = h0._ = void 0;
  var r = M9(), K1 = oY(), Q6 = M9();
  Object.defineProperty(h0, "_", { enumerable: true, get: function() {
    return Q6._;
  } });
  Object.defineProperty(h0, "str", { enumerable: true, get: function() {
    return Q6.str;
  } });
  Object.defineProperty(h0, "strConcat", { enumerable: true, get: function() {
    return Q6.strConcat;
  } });
  Object.defineProperty(h0, "nil", { enumerable: true, get: function() {
    return Q6.nil;
  } });
  Object.defineProperty(h0, "getProperty", { enumerable: true, get: function() {
    return Q6.getProperty;
  } });
  Object.defineProperty(h0, "stringify", { enumerable: true, get: function() {
    return Q6.stringify;
  } });
  Object.defineProperty(h0, "regexpCode", { enumerable: true, get: function() {
    return Q6.regexpCode;
  } });
  Object.defineProperty(h0, "Name", { enumerable: true, get: function() {
    return Q6.Name;
  } });
  var BQ = oY();
  Object.defineProperty(h0, "Scope", { enumerable: true, get: function() {
    return BQ.Scope;
  } });
  Object.defineProperty(h0, "ValueScope", { enumerable: true, get: function() {
    return BQ.ValueScope;
  } });
  Object.defineProperty(h0, "ValueScopeName", { enumerable: true, get: function() {
    return BQ.ValueScopeName;
  } });
  Object.defineProperty(h0, "varKinds", { enumerable: true, get: function() {
    return BQ.varKinds;
  } });
  h0.operators = { GT: new r._Code(">"), GTE: new r._Code(">="), LT: new r._Code("<"), LTE: new r._Code("<="), EQ: new r._Code("==="), NEQ: new r._Code("!=="), NOT: new r._Code("!"), OR: new r._Code("||"), AND: new r._Code("&&"), ADD: new r._Code("+") };
  class $6 {
    optimizeNodes() {
      return this;
    }
    optimizeNames(Q, $) {
      return this;
    }
  }
  class w3 extends $6 {
    constructor(Q, $, X) {
      super();
      this.varKind = Q, this.name = $, this.rhs = X;
    }
    render({ es5: Q, _n: $ }) {
      let X = Q ? K1.varKinds.var : this.varKind, Y = this.rhs === void 0 ? "" : ` = ${this.rhs}`;
      return `${X} ${this.name}${Y};` + $;
    }
    optimizeNames(Q, $) {
      if (!Q[this.name.str]) return;
      if (this.rhs) this.rhs = A4(this.rhs, Q, $);
      return this;
    }
    get names() {
      return this.rhs instanceof r._CodeOrName ? this.rhs.names : {};
    }
  }
  class aY extends $6 {
    constructor(Q, $, X) {
      super();
      this.lhs = Q, this.rhs = $, this.sideEffects = X;
    }
    render({ _n: Q }) {
      return `${this.lhs} = ${this.rhs};` + Q;
    }
    optimizeNames(Q, $) {
      if (this.lhs instanceof r.Name && !Q[this.lhs.str] && !this.sideEffects) return;
      return this.rhs = A4(this.rhs, Q, $), this;
    }
    get names() {
      let Q = this.lhs instanceof r.Name ? {} : { ...this.lhs.names };
      return HQ(Q, this.rhs);
    }
  }
  class M3 extends aY {
    constructor(Q, $, X, Y) {
      super(Q, X, Y);
      this.op = $;
    }
    render({ _n: Q }) {
      return `${this.lhs} ${this.op}= ${this.rhs};` + Q;
    }
  }
  class A3 extends $6 {
    constructor(Q) {
      super();
      this.label = Q, this.names = {};
    }
    render({ _n: Q }) {
      return `${this.label}:` + Q;
    }
  }
  class j3 extends $6 {
    constructor(Q) {
      super();
      this.label = Q, this.names = {};
    }
    render({ _n: Q }) {
      return `break${this.label ? ` ${this.label}` : ""};` + Q;
    }
  }
  class R3 extends $6 {
    constructor(Q) {
      super();
      this.error = Q;
    }
    render({ _n: Q }) {
      return `throw ${this.error};` + Q;
    }
    get names() {
      return this.error.names;
    }
  }
  class I3 extends $6 {
    constructor(Q) {
      super();
      this.code = Q;
    }
    render({ _n: Q }) {
      return `${this.code};` + Q;
    }
    optimizeNodes() {
      return `${this.code}` ? this : void 0;
    }
    optimizeNames(Q, $) {
      return this.code = A4(this.code, Q, $), this;
    }
    get names() {
      return this.code instanceof r._CodeOrName ? this.code.names : {};
    }
  }
  class zQ extends $6 {
    constructor(Q = []) {
      super();
      this.nodes = Q;
    }
    render(Q) {
      return this.nodes.reduce(($, X) => $ + X.render(Q), "");
    }
    optimizeNodes() {
      let { nodes: Q } = this, $ = Q.length;
      while ($--) {
        let X = Q[$].optimizeNodes();
        if (Array.isArray(X)) Q.splice($, 1, ...X);
        else if (X) Q[$] = X;
        else Q.splice($, 1);
      }
      return Q.length > 0 ? this : void 0;
    }
    optimizeNames(Q, $) {
      let { nodes: X } = this, Y = X.length;
      while (Y--) {
        let J = X[Y];
        if (J.optimizeNames(Q, $)) continue;
        nw(Q, J.names), X.splice(Y, 1);
      }
      return X.length > 0 ? this : void 0;
    }
    get names() {
      return this.nodes.reduce((Q, $) => b6(Q, $.names), {});
    }
  }
  class X6 extends zQ {
    render(Q) {
      return "{" + Q._n + super.render(Q) + "}" + Q._n;
    }
  }
  class E3 extends zQ {
  }
  class A9 extends X6 {
  }
  A9.kind = "else";
  class _1 extends X6 {
    constructor(Q, $) {
      super($);
      this.condition = Q;
    }
    render(Q) {
      let $ = `if(${this.condition})` + super.render(Q);
      if (this.else) $ += "else " + this.else.render(Q);
      return $;
    }
    optimizeNodes() {
      super.optimizeNodes();
      let Q = this.condition;
      if (Q === true) return this.nodes;
      let $ = this.else;
      if ($) {
        let X = $.optimizeNodes();
        $ = this.else = Array.isArray(X) ? new A9(X) : X;
      }
      if ($) {
        if (Q === false) return $ instanceof _1 ? $ : $.nodes;
        if (this.nodes.length) return this;
        return new _1(S3(Q), $ instanceof _1 ? [$] : $.nodes);
      }
      if (Q === false || !this.nodes.length) return;
      return this;
    }
    optimizeNames(Q, $) {
      var X;
      if (this.else = (X = this.else) === null || X === void 0 ? void 0 : X.optimizeNames(Q, $), !(super.optimizeNames(Q, $) || this.else)) return;
      return this.condition = A4(this.condition, Q, $), this;
    }
    get names() {
      let Q = super.names;
      if (HQ(Q, this.condition), this.else) b6(Q, this.else.names);
      return Q;
    }
  }
  _1.kind = "if";
  class M4 extends X6 {
  }
  M4.kind = "for";
  class P3 extends M4 {
    constructor(Q) {
      super();
      this.iteration = Q;
    }
    render(Q) {
      return `for(${this.iteration})` + super.render(Q);
    }
    optimizeNames(Q, $) {
      if (!super.optimizeNames(Q, $)) return;
      return this.iteration = A4(this.iteration, Q, $), this;
    }
    get names() {
      return b6(super.names, this.iteration.names);
    }
  }
  class b3 extends M4 {
    constructor(Q, $, X, Y) {
      super();
      this.varKind = Q, this.name = $, this.from = X, this.to = Y;
    }
    render(Q) {
      let $ = Q.es5 ? K1.varKinds.var : this.varKind, { name: X, from: Y, to: J } = this;
      return `for(${$} ${X}=${Y}; ${X}<${J}; ${X}++)` + super.render(Q);
    }
    get names() {
      let Q = HQ(super.names, this.from);
      return HQ(Q, this.to);
    }
  }
  class rY extends M4 {
    constructor(Q, $, X, Y) {
      super();
      this.loop = Q, this.varKind = $, this.name = X, this.iterable = Y;
    }
    render(Q) {
      return `for(${this.varKind} ${this.name} ${this.loop} ${this.iterable})` + super.render(Q);
    }
    optimizeNames(Q, $) {
      if (!super.optimizeNames(Q, $)) return;
      return this.iterable = A4(this.iterable, Q, $), this;
    }
    get names() {
      return b6(super.names, this.iterable.names);
    }
  }
  class YQ extends X6 {
    constructor(Q, $, X) {
      super();
      this.name = Q, this.args = $, this.async = X;
    }
    render(Q) {
      return `${this.async ? "async " : ""}function ${this.name}(${this.args})` + super.render(Q);
    }
  }
  YQ.kind = "func";
  class JQ extends zQ {
    render(Q) {
      return "return " + super.render(Q);
    }
  }
  JQ.kind = "return";
  class Z3 extends X6 {
    render(Q) {
      let $ = "try" + super.render(Q);
      if (this.catch) $ += this.catch.render(Q);
      if (this.finally) $ += this.finally.render(Q);
      return $;
    }
    optimizeNodes() {
      var Q, $;
      return super.optimizeNodes(), (Q = this.catch) === null || Q === void 0 || Q.optimizeNodes(), ($ = this.finally) === null || $ === void 0 || $.optimizeNodes(), this;
    }
    optimizeNames(Q, $) {
      var X, Y;
      return super.optimizeNames(Q, $), (X = this.catch) === null || X === void 0 || X.optimizeNames(Q, $), (Y = this.finally) === null || Y === void 0 || Y.optimizeNames(Q, $), this;
    }
    get names() {
      let Q = super.names;
      if (this.catch) b6(Q, this.catch.names);
      if (this.finally) b6(Q, this.finally.names);
      return Q;
    }
  }
  class GQ extends X6 {
    constructor(Q) {
      super();
      this.error = Q;
    }
    render(Q) {
      return `catch(${this.error})` + super.render(Q);
    }
  }
  GQ.kind = "catch";
  class WQ extends X6 {
    render(Q) {
      return "finally" + super.render(Q);
    }
  }
  WQ.kind = "finally";
  class C3 {
    constructor(Q, $ = {}) {
      this._values = {}, this._blockStarts = [], this._constants = {}, this.opts = { ...$, _n: $.lines ? `
` : "" }, this._extScope = Q, this._scope = new K1.Scope({ parent: Q }), this._nodes = [new E3()];
    }
    toString() {
      return this._root.render(this.opts);
    }
    name(Q) {
      return this._scope.name(Q);
    }
    scopeName(Q) {
      return this._extScope.name(Q);
    }
    scopeValue(Q, $) {
      let X = this._extScope.value(Q, $);
      return (this._values[X.prefix] || (this._values[X.prefix] = /* @__PURE__ */ new Set())).add(X), X;
    }
    getScopeValue(Q, $) {
      return this._extScope.getValue(Q, $);
    }
    scopeRefs(Q) {
      return this._extScope.scopeRefs(Q, this._values);
    }
    scopeCode() {
      return this._extScope.scopeCode(this._values);
    }
    _def(Q, $, X, Y) {
      let J = this._scope.toName($);
      if (X !== void 0 && Y) this._constants[J.str] = X;
      return this._leafNode(new w3(Q, J, X)), J;
    }
    const(Q, $, X) {
      return this._def(K1.varKinds.const, Q, $, X);
    }
    let(Q, $, X) {
      return this._def(K1.varKinds.let, Q, $, X);
    }
    var(Q, $, X) {
      return this._def(K1.varKinds.var, Q, $, X);
    }
    assign(Q, $, X) {
      return this._leafNode(new aY(Q, $, X));
    }
    add(Q, $) {
      return this._leafNode(new M3(Q, h0.operators.ADD, $));
    }
    code(Q) {
      if (typeof Q == "function") Q();
      else if (Q !== r.nil) this._leafNode(new I3(Q));
      return this;
    }
    object(...Q) {
      let $ = ["{"];
      for (let [X, Y] of Q) {
        if ($.length > 1) $.push(",");
        if ($.push(X), X !== Y || this.opts.es5) $.push(":"), (0, r.addCodeArg)($, Y);
      }
      return $.push("}"), new r._Code($);
    }
    if(Q, $, X) {
      if (this._blockNode(new _1(Q)), $ && X) this.code($).else().code(X).endIf();
      else if ($) this.code($).endIf();
      else if (X) throw Error('CodeGen: "else" body without "then" body');
      return this;
    }
    elseIf(Q) {
      return this._elseNode(new _1(Q));
    }
    else() {
      return this._elseNode(new A9());
    }
    endIf() {
      return this._endBlockNode(_1, A9);
    }
    _for(Q, $) {
      if (this._blockNode(Q), $) this.code($).endFor();
      return this;
    }
    for(Q, $) {
      return this._for(new P3(Q), $);
    }
    forRange(Q, $, X, Y, J = this.opts.es5 ? K1.varKinds.var : K1.varKinds.let) {
      let G = this._scope.toName(Q);
      return this._for(new b3(J, G, $, X), () => Y(G));
    }
    forOf(Q, $, X, Y = K1.varKinds.const) {
      let J = this._scope.toName(Q);
      if (this.opts.es5) {
        let G = $ instanceof r.Name ? $ : this.var("_arr", $);
        return this.forRange("_i", 0, r._`${G}.length`, (W) => {
          this.var(J, r._`${G}[${W}]`), X(J);
        });
      }
      return this._for(new rY("of", Y, J, $), () => X(J));
    }
    forIn(Q, $, X, Y = this.opts.es5 ? K1.varKinds.var : K1.varKinds.const) {
      if (this.opts.ownProperties) return this.forOf(Q, r._`Object.keys(${$})`, X);
      let J = this._scope.toName(Q);
      return this._for(new rY("in", Y, J, $), () => X(J));
    }
    endFor() {
      return this._endBlockNode(M4);
    }
    label(Q) {
      return this._leafNode(new A3(Q));
    }
    break(Q) {
      return this._leafNode(new j3(Q));
    }
    return(Q) {
      let $ = new JQ();
      if (this._blockNode($), this.code(Q), $.nodes.length !== 1) throw Error('CodeGen: "return" should have one node');
      return this._endBlockNode(JQ);
    }
    try(Q, $, X) {
      if (!$ && !X) throw Error('CodeGen: "try" without "catch" and "finally"');
      let Y = new Z3();
      if (this._blockNode(Y), this.code(Q), $) {
        let J = this.name("e");
        this._currNode = Y.catch = new GQ(J), $(J);
      }
      if (X) this._currNode = Y.finally = new WQ(), this.code(X);
      return this._endBlockNode(GQ, WQ);
    }
    throw(Q) {
      return this._leafNode(new R3(Q));
    }
    block(Q, $) {
      if (this._blockStarts.push(this._nodes.length), Q) this.code(Q).endBlock($);
      return this;
    }
    endBlock(Q) {
      let $ = this._blockStarts.pop();
      if ($ === void 0) throw Error("CodeGen: not in self-balancing block");
      let X = this._nodes.length - $;
      if (X < 0 || Q !== void 0 && X !== Q) throw Error(`CodeGen: wrong number of nodes: ${X} vs ${Q} expected`);
      return this._nodes.length = $, this;
    }
    func(Q, $ = r.nil, X, Y) {
      if (this._blockNode(new YQ(Q, $, X)), Y) this.code(Y).endFunc();
      return this;
    }
    endFunc() {
      return this._endBlockNode(YQ);
    }
    optimize(Q = 1) {
      while (Q-- > 0) this._root.optimizeNodes(), this._root.optimizeNames(this._root.names, this._constants);
    }
    _leafNode(Q) {
      return this._currNode.nodes.push(Q), this;
    }
    _blockNode(Q) {
      this._currNode.nodes.push(Q), this._nodes.push(Q);
    }
    _endBlockNode(Q, $) {
      let X = this._currNode;
      if (X instanceof Q || $ && X instanceof $) return this._nodes.pop(), this;
      throw Error(`CodeGen: not in block "${$ ? `${Q.kind}/${$.kind}` : Q.kind}"`);
    }
    _elseNode(Q) {
      let $ = this._currNode;
      if (!($ instanceof _1)) throw Error('CodeGen: "else" without "if"');
      return this._currNode = $.else = Q, this;
    }
    get _root() {
      return this._nodes[0];
    }
    get _currNode() {
      let Q = this._nodes;
      return Q[Q.length - 1];
    }
    set _currNode(Q) {
      let $ = this._nodes;
      $[$.length - 1] = Q;
    }
  }
  h0.CodeGen = C3;
  function b6(Q, $) {
    for (let X in $) Q[X] = (Q[X] || 0) + ($[X] || 0);
    return Q;
  }
  function HQ(Q, $) {
    return $ instanceof r._CodeOrName ? b6(Q, $.names) : Q;
  }
  function A4(Q, $, X) {
    if (Q instanceof r.Name) return Y(Q);
    if (!J(Q)) return Q;
    return new r._Code(Q._items.reduce((G, W) => {
      if (W instanceof r.Name) W = Y(W);
      if (W instanceof r._Code) G.push(...W._items);
      else G.push(W);
      return G;
    }, []));
    function Y(G) {
      let W = X[G.str];
      if (W === void 0 || $[G.str] !== 1) return G;
      return delete $[G.str], W;
    }
    function J(G) {
      return G instanceof r._Code && G._items.some((W) => W instanceof r.Name && $[W.str] === 1 && X[W.str] !== void 0);
    }
  }
  function nw(Q, $) {
    for (let X in $) Q[X] = (Q[X] || 0) - ($[X] || 0);
  }
  function S3(Q) {
    return typeof Q == "boolean" || typeof Q == "number" || Q === null ? !Q : r._`!${tY(Q)}`;
  }
  h0.not = S3;
  var ow = _3(h0.operators.AND);
  function rw(...Q) {
    return Q.reduce(ow);
  }
  h0.and = rw;
  var tw = _3(h0.operators.OR);
  function aw(...Q) {
    return Q.reduce(tw);
  }
  h0.or = aw;
  function _3(Q) {
    return ($, X) => $ === r.nil ? X : X === r.nil ? $ : r._`${tY($)} ${Q} ${tY(X)}`;
  }
  function tY(Q) {
    return Q instanceof r.Name ? Q : r._`(${Q})`;
  }
});
var a = P((u3) => {
  Object.defineProperty(u3, "__esModule", { value: true });
  u3.checkStrictMode = u3.getErrorPath = u3.Type = u3.useFunc = u3.setEvaluated = u3.evaluatedPropsToName = u3.mergeEvaluated = u3.eachItem = u3.unescapeJsonPointer = u3.escapeJsonPointer = u3.escapeFragment = u3.unescapeFragment = u3.schemaRefOrVal = u3.schemaHasRulesButRef = u3.schemaHasRules = u3.checkUnknownRules = u3.alwaysValidSchema = u3.toHash = void 0;
  var Q0 = c(), $M = M9();
  function XM(Q) {
    let $ = {};
    for (let X of Q) $[X] = true;
    return $;
  }
  u3.toHash = XM;
  function YM(Q, $) {
    if (typeof $ == "boolean") return $;
    if (Object.keys($).length === 0) return true;
    return x3(Q, $), !y3($, Q.self.RULES.all);
  }
  u3.alwaysValidSchema = YM;
  function x3(Q, $ = Q.schema) {
    let { opts: X, self: Y } = Q;
    if (!X.strictSchema) return;
    if (typeof $ === "boolean") return;
    let J = Y.RULES.keywords;
    for (let G in $) if (!J[G]) f3(Q, `unknown keyword: "${G}"`);
  }
  u3.checkUnknownRules = x3;
  function y3(Q, $) {
    if (typeof Q == "boolean") return !Q;
    for (let X in Q) if ($[X]) return true;
    return false;
  }
  u3.schemaHasRules = y3;
  function JM(Q, $) {
    if (typeof Q == "boolean") return !Q;
    for (let X in Q) if (X !== "$ref" && $.all[X]) return true;
    return false;
  }
  u3.schemaHasRulesButRef = JM;
  function GM({ topSchemaRef: Q, schemaPath: $ }, X, Y, J) {
    if (!J) {
      if (typeof X == "number" || typeof X == "boolean") return X;
      if (typeof X == "string") return Q0._`${X}`;
    }
    return Q0._`${Q}${$}${(0, Q0.getProperty)(Y)}`;
  }
  u3.schemaRefOrVal = GM;
  function WM(Q) {
    return g3(decodeURIComponent(Q));
  }
  u3.unescapeFragment = WM;
  function HM(Q) {
    return encodeURIComponent(eY(Q));
  }
  u3.escapeFragment = HM;
  function eY(Q) {
    if (typeof Q == "number") return `${Q}`;
    return Q.replace(/~/g, "~0").replace(/\//g, "~1");
  }
  u3.escapeJsonPointer = eY;
  function g3(Q) {
    return Q.replace(/~1/g, "/").replace(/~0/g, "~");
  }
  u3.unescapeJsonPointer = g3;
  function BM(Q, $) {
    if (Array.isArray(Q)) for (let X of Q) $(X);
    else $(Q);
  }
  u3.eachItem = BM;
  function v3({ mergeNames: Q, mergeToName: $, mergeValues: X, resultToName: Y }) {
    return (J, G, W, H) => {
      let B = W === void 0 ? G : W instanceof Q0.Name ? (G instanceof Q0.Name ? Q(J, G, W) : $(J, G, W), W) : G instanceof Q0.Name ? ($(J, W, G), G) : X(G, W);
      return H === Q0.Name && !(B instanceof Q0.Name) ? Y(J, B) : B;
    };
  }
  u3.mergeEvaluated = { props: v3({ mergeNames: (Q, $, X) => Q.if(Q0._`${X} !== true && ${$} !== undefined`, () => {
    Q.if(Q0._`${$} === true`, () => Q.assign(X, true), () => Q.assign(X, Q0._`${X} || {}`).code(Q0._`Object.assign(${X}, ${$})`));
  }), mergeToName: (Q, $, X) => Q.if(Q0._`${X} !== true`, () => {
    if ($ === true) Q.assign(X, true);
    else Q.assign(X, Q0._`${X} || {}`), Q7(Q, X, $);
  }), mergeValues: (Q, $) => Q === true ? true : { ...Q, ...$ }, resultToName: h3 }), items: v3({ mergeNames: (Q, $, X) => Q.if(Q0._`${X} !== true && ${$} !== undefined`, () => Q.assign(X, Q0._`${$} === true ? true : ${X} > ${$} ? ${X} : ${$}`)), mergeToName: (Q, $, X) => Q.if(Q0._`${X} !== true`, () => Q.assign(X, $ === true ? true : Q0._`${X} > ${$} ? ${X} : ${$}`)), mergeValues: (Q, $) => Q === true ? true : Math.max(Q, $), resultToName: (Q, $) => Q.var("items", $) }) };
  function h3(Q, $) {
    if ($ === true) return Q.var("props", true);
    let X = Q.var("props", Q0._`{}`);
    if ($ !== void 0) Q7(Q, X, $);
    return X;
  }
  u3.evaluatedPropsToName = h3;
  function Q7(Q, $, X) {
    Object.keys(X).forEach((Y) => Q.assign(Q0._`${$}${(0, Q0.getProperty)(Y)}`, true));
  }
  u3.setEvaluated = Q7;
  var T3 = {};
  function zM(Q, $) {
    return Q.scopeValue("func", { ref: $, code: T3[$.code] || (T3[$.code] = new $M._Code($.code)) });
  }
  u3.useFunc = zM;
  var sY;
  (function(Q) {
    Q[Q.Num = 0] = "Num", Q[Q.Str = 1] = "Str";
  })(sY || (u3.Type = sY = {}));
  function KM(Q, $, X) {
    if (Q instanceof Q0.Name) {
      let Y = $ === sY.Num;
      return X ? Y ? Q0._`"[" + ${Q} + "]"` : Q0._`"['" + ${Q} + "']"` : Y ? Q0._`"/" + ${Q}` : Q0._`"/" + ${Q}.replace(/~/g, "~0").replace(/\\//g, "~1")`;
    }
    return X ? (0, Q0.getProperty)(Q).toString() : "/" + eY(Q);
  }
  u3.getErrorPath = KM;
  function f3(Q, $, X = Q.opts.strictSchema) {
    if (!X) return;
    if ($ = `strict mode: ${$}`, X === true) throw Error($);
    Q.self.logger.warn($);
  }
  u3.checkStrictMode = f3;
});
var k1 = P((l3) => {
  Object.defineProperty(l3, "__esModule", { value: true });
  var Z0 = c(), ZM = { data: new Z0.Name("data"), valCxt: new Z0.Name("valCxt"), instancePath: new Z0.Name("instancePath"), parentData: new Z0.Name("parentData"), parentDataProperty: new Z0.Name("parentDataProperty"), rootData: new Z0.Name("rootData"), dynamicAnchors: new Z0.Name("dynamicAnchors"), vErrors: new Z0.Name("vErrors"), errors: new Z0.Name("errors"), this: new Z0.Name("this"), self: new Z0.Name("self"), scope: new Z0.Name("scope"), json: new Z0.Name("json"), jsonPos: new Z0.Name("jsonPos"), jsonLen: new Z0.Name("jsonLen"), jsonPart: new Z0.Name("jsonPart") };
  l3.default = ZM;
});
var j9 = P((i3) => {
  Object.defineProperty(i3, "__esModule", { value: true });
  i3.extendErrors = i3.resetErrorsCount = i3.reportExtraError = i3.reportError = i3.keyword$DataError = i3.keywordError = void 0;
  var t = c(), VQ = a(), k0 = k1();
  i3.keywordError = { message: ({ keyword: Q }) => t.str`must pass "${Q}" keyword validation` };
  i3.keyword$DataError = { message: ({ keyword: Q, schemaType: $ }) => $ ? t.str`"${Q}" keyword must be ${$} ($data)` : t.str`"${Q}" keyword is invalid ($data)` };
  function SM(Q, $ = i3.keywordError, X, Y) {
    let { it: J } = Q, { gen: G, compositeRule: W, allErrors: H } = J, B = d3(Q, $, X);
    if (Y !== null && Y !== void 0 ? Y : W || H) c3(G, B);
    else p3(J, t._`[${B}]`);
  }
  i3.reportError = SM;
  function _M(Q, $ = i3.keywordError, X) {
    let { it: Y } = Q, { gen: J, compositeRule: G, allErrors: W } = Y, H = d3(Q, $, X);
    if (c3(J, H), !(G || W)) p3(Y, k0.default.vErrors);
  }
  i3.reportExtraError = _M;
  function kM(Q, $) {
    Q.assign(k0.default.errors, $), Q.if(t._`${k0.default.vErrors} !== null`, () => Q.if($, () => Q.assign(t._`${k0.default.vErrors}.length`, $), () => Q.assign(k0.default.vErrors, null)));
  }
  i3.resetErrorsCount = kM;
  function vM({ gen: Q, keyword: $, schemaValue: X, data: Y, errsCount: J, it: G }) {
    if (J === void 0) throw Error("ajv implementation error");
    let W = Q.name("err");
    Q.forRange("i", J, k0.default.errors, (H) => {
      if (Q.const(W, t._`${k0.default.vErrors}[${H}]`), Q.if(t._`${W}.instancePath === undefined`, () => Q.assign(t._`${W}.instancePath`, (0, t.strConcat)(k0.default.instancePath, G.errorPath))), Q.assign(t._`${W}.schemaPath`, t.str`${G.errSchemaPath}/${$}`), G.opts.verbose) Q.assign(t._`${W}.schema`, X), Q.assign(t._`${W}.data`, Y);
    });
  }
  i3.extendErrors = vM;
  function c3(Q, $) {
    let X = Q.const("err", $);
    Q.if(t._`${k0.default.vErrors} === null`, () => Q.assign(k0.default.vErrors, t._`[${X}]`), t._`${k0.default.vErrors}.push(${X})`), Q.code(t._`${k0.default.errors}++`);
  }
  function p3(Q, $) {
    let { gen: X, validateName: Y, schemaEnv: J } = Q;
    if (J.$async) X.throw(t._`new ${Q.ValidationError}(${$})`);
    else X.assign(t._`${Y}.errors`, $), X.return(false);
  }
  var Z6 = { keyword: new t.Name("keyword"), schemaPath: new t.Name("schemaPath"), params: new t.Name("params"), propertyName: new t.Name("propertyName"), message: new t.Name("message"), schema: new t.Name("schema"), parentSchema: new t.Name("parentSchema") };
  function d3(Q, $, X) {
    let { createErrors: Y } = Q.it;
    if (Y === false) return t._`{}`;
    return TM(Q, $, X);
  }
  function TM(Q, $, X = {}) {
    let { gen: Y, it: J } = Q, G = [xM(J, X), yM(Q, X)];
    return gM(Q, $, G), Y.object(...G);
  }
  function xM({ errorPath: Q }, { instancePath: $ }) {
    let X = $ ? t.str`${Q}${(0, VQ.getErrorPath)($, VQ.Type.Str)}` : Q;
    return [k0.default.instancePath, (0, t.strConcat)(k0.default.instancePath, X)];
  }
  function yM({ keyword: Q, it: { errSchemaPath: $ } }, { schemaPath: X, parentSchema: Y }) {
    let J = Y ? $ : t.str`${$}/${Q}`;
    if (X) J = t.str`${J}${(0, VQ.getErrorPath)(X, VQ.Type.Str)}`;
    return [Z6.schemaPath, J];
  }
  function gM(Q, { params: $, message: X }, Y) {
    let { keyword: J, data: G, schemaValue: W, it: H } = Q, { opts: B, propertyName: z, topSchemaRef: K, schemaPath: U } = H;
    if (Y.push([Z6.keyword, J], [Z6.params, typeof $ == "function" ? $(Q) : $ || t._`{}`]), B.messages) Y.push([Z6.message, typeof X == "function" ? X(Q) : X]);
    if (B.verbose) Y.push([Z6.schema, W], [Z6.parentSchema, t._`${K}${U}`], [k0.default.data, G]);
    if (z) Y.push([Z6.propertyName, z]);
  }
});
var a3 = P((r3) => {
  Object.defineProperty(r3, "__esModule", { value: true });
  r3.boolOrEmptySchema = r3.topBoolOrEmptySchema = void 0;
  var lM = j9(), cM = c(), pM = k1(), dM = { message: "boolean schema is false" };
  function iM(Q) {
    let { gen: $, schema: X, validateName: Y } = Q;
    if (X === false) o3(Q, false);
    else if (typeof X == "object" && X.$async === true) $.return(pM.default.data);
    else $.assign(cM._`${Y}.errors`, null), $.return(true);
  }
  r3.topBoolOrEmptySchema = iM;
  function nM(Q, $) {
    let { gen: X, schema: Y } = Q;
    if (Y === false) X.var($, false), o3(Q);
    else X.var($, true);
  }
  r3.boolOrEmptySchema = nM;
  function o3(Q, $) {
    let { gen: X, data: Y } = Q, J = { gen: X, keyword: "false schema", data: Y, schema: false, schemaCode: false, schemaValue: false, params: {}, it: Q };
    (0, lM.reportError)(J, dM, void 0, $);
  }
});
var X7 = P((s3) => {
  Object.defineProperty(s3, "__esModule", { value: true });
  s3.getRules = s3.isJSONType = void 0;
  var rM = ["string", "number", "integer", "boolean", "null", "object", "array"], tM = new Set(rM);
  function aM(Q) {
    return typeof Q == "string" && tM.has(Q);
  }
  s3.isJSONType = aM;
  function sM() {
    let Q = { number: { type: "number", rules: [] }, string: { type: "string", rules: [] }, array: { type: "array", rules: [] }, object: { type: "object", rules: [] } };
    return { types: { ...Q, integer: true, boolean: true, null: true }, rules: [{ rules: [] }, Q.number, Q.string, Q.array, Q.object], post: { rules: [] }, all: {}, keywords: {} };
  }
  s3.getRules = sM;
});
var Y7 = P((XH) => {
  Object.defineProperty(XH, "__esModule", { value: true });
  XH.shouldUseRule = XH.shouldUseGroup = XH.schemaHasRulesForType = void 0;
  function QA({ schema: Q, self: $ }, X) {
    let Y = $.RULES.types[X];
    return Y && Y !== true && QH(Q, Y);
  }
  XH.schemaHasRulesForType = QA;
  function QH(Q, $) {
    return $.rules.some((X) => $H(Q, X));
  }
  XH.shouldUseGroup = QH;
  function $H(Q, $) {
    var X;
    return Q[$.keyword] !== void 0 || ((X = $.definition.implements) === null || X === void 0 ? void 0 : X.some((Y) => Q[Y] !== void 0));
  }
  XH.shouldUseRule = $H;
});
var R9 = P((HH) => {
  Object.defineProperty(HH, "__esModule", { value: true });
  HH.reportTypeError = HH.checkDataTypes = HH.checkDataType = HH.coerceAndCheckDataType = HH.getJSONTypes = HH.getSchemaTypes = HH.DataType = void 0;
  var YA = X7(), JA = Y7(), GA = j9(), l = c(), JH = a(), j4;
  (function(Q) {
    Q[Q.Correct = 0] = "Correct", Q[Q.Wrong = 1] = "Wrong";
  })(j4 || (HH.DataType = j4 = {}));
  function WA(Q) {
    let $ = GH(Q.type);
    if ($.includes("null")) {
      if (Q.nullable === false) throw Error("type: null contradicts nullable: false");
    } else {
      if (!$.length && Q.nullable !== void 0) throw Error('"nullable" cannot be used without "type"');
      if (Q.nullable === true) $.push("null");
    }
    return $;
  }
  HH.getSchemaTypes = WA;
  function GH(Q) {
    let $ = Array.isArray(Q) ? Q : Q ? [Q] : [];
    if ($.every(YA.isJSONType)) return $;
    throw Error("type must be JSONType or JSONType[]: " + $.join(","));
  }
  HH.getJSONTypes = GH;
  function HA(Q, $) {
    let { gen: X, data: Y, opts: J } = Q, G = BA($, J.coerceTypes), W = $.length > 0 && !(G.length === 0 && $.length === 1 && (0, JA.schemaHasRulesForType)(Q, $[0]));
    if (W) {
      let H = G7($, Y, J.strictNumbers, j4.Wrong);
      X.if(H, () => {
        if (G.length) zA(Q, $, G);
        else W7(Q);
      });
    }
    return W;
  }
  HH.coerceAndCheckDataType = HA;
  var WH = /* @__PURE__ */ new Set(["string", "number", "integer", "boolean", "null"]);
  function BA(Q, $) {
    return $ ? Q.filter((X) => WH.has(X) || $ === "array" && X === "array") : [];
  }
  function zA(Q, $, X) {
    let { gen: Y, data: J, opts: G } = Q, W = Y.let("dataType", l._`typeof ${J}`), H = Y.let("coerced", l._`undefined`);
    if (G.coerceTypes === "array") Y.if(l._`${W} == 'object' && Array.isArray(${J}) && ${J}.length == 1`, () => Y.assign(J, l._`${J}[0]`).assign(W, l._`typeof ${J}`).if(G7($, J, G.strictNumbers), () => Y.assign(H, J)));
    Y.if(l._`${H} !== undefined`);
    for (let z of X) if (WH.has(z) || z === "array" && G.coerceTypes === "array") B(z);
    Y.else(), W7(Q), Y.endIf(), Y.if(l._`${H} !== undefined`, () => {
      Y.assign(J, H), KA(Q, H);
    });
    function B(z) {
      switch (z) {
        case "string":
          Y.elseIf(l._`${W} == "number" || ${W} == "boolean"`).assign(H, l._`"" + ${J}`).elseIf(l._`${J} === null`).assign(H, l._`""`);
          return;
        case "number":
          Y.elseIf(l._`${W} == "boolean" || ${J} === null
              || (${W} == "string" && ${J} && ${J} == +${J})`).assign(H, l._`+${J}`);
          return;
        case "integer":
          Y.elseIf(l._`${W} === "boolean" || ${J} === null
              || (${W} === "string" && ${J} && ${J} == +${J} && !(${J} % 1))`).assign(H, l._`+${J}`);
          return;
        case "boolean":
          Y.elseIf(l._`${J} === "false" || ${J} === 0 || ${J} === null`).assign(H, false).elseIf(l._`${J} === "true" || ${J} === 1`).assign(H, true);
          return;
        case "null":
          Y.elseIf(l._`${J} === "" || ${J} === 0 || ${J} === false`), Y.assign(H, null);
          return;
        case "array":
          Y.elseIf(l._`${W} === "string" || ${W} === "number"
              || ${W} === "boolean" || ${J} === null`).assign(H, l._`[${J}]`);
      }
    }
  }
  function KA({ gen: Q, parentData: $, parentDataProperty: X }, Y) {
    Q.if(l._`${$} !== undefined`, () => Q.assign(l._`${$}[${X}]`, Y));
  }
  function J7(Q, $, X, Y = j4.Correct) {
    let J = Y === j4.Correct ? l.operators.EQ : l.operators.NEQ, G;
    switch (Q) {
      case "null":
        return l._`${$} ${J} null`;
      case "array":
        G = l._`Array.isArray(${$})`;
        break;
      case "object":
        G = l._`${$} && typeof ${$} == "object" && !Array.isArray(${$})`;
        break;
      case "integer":
        G = W(l._`!(${$} % 1) && !isNaN(${$})`);
        break;
      case "number":
        G = W();
        break;
      default:
        return l._`typeof ${$} ${J} ${Q}`;
    }
    return Y === j4.Correct ? G : (0, l.not)(G);
    function W(H = l.nil) {
      return (0, l.and)(l._`typeof ${$} == "number"`, H, X ? l._`isFinite(${$})` : l.nil);
    }
  }
  HH.checkDataType = J7;
  function G7(Q, $, X, Y) {
    if (Q.length === 1) return J7(Q[0], $, X, Y);
    let J, G = (0, JH.toHash)(Q);
    if (G.array && G.object) {
      let W = l._`typeof ${$} != "object"`;
      J = G.null ? W : l._`!${$} || ${W}`, delete G.null, delete G.array, delete G.object;
    } else J = l.nil;
    if (G.number) delete G.integer;
    for (let W in G) J = (0, l.and)(J, J7(W, $, X, Y));
    return J;
  }
  HH.checkDataTypes = G7;
  var VA = { message: ({ schema: Q }) => `must be ${Q}`, params: ({ schema: Q, schemaValue: $ }) => typeof Q == "string" ? l._`{type: ${Q}}` : l._`{type: ${$}}` };
  function W7(Q) {
    let $ = qA(Q);
    (0, GA.reportError)($, VA);
  }
  HH.reportTypeError = W7;
  function qA(Q) {
    let { gen: $, data: X, schema: Y } = Q, J = (0, JH.schemaRefOrVal)(Q, Y, "type");
    return { gen: $, keyword: "type", data: X, schema: Y.type, schemaCode: J, schemaValue: J, parentSchema: Y, params: {}, it: Q };
  }
});
var qH = P((KH) => {
  Object.defineProperty(KH, "__esModule", { value: true });
  KH.assignDefaults = void 0;
  var R4 = c(), wA = a();
  function MA(Q, $) {
    let { properties: X, items: Y } = Q.schema;
    if ($ === "object" && X) for (let J in X) zH(Q, J, X[J].default);
    else if ($ === "array" && Array.isArray(Y)) Y.forEach((J, G) => zH(Q, G, J.default));
  }
  KH.assignDefaults = MA;
  function zH(Q, $, X) {
    let { gen: Y, compositeRule: J, data: G, opts: W } = Q;
    if (X === void 0) return;
    let H = R4._`${G}${(0, R4.getProperty)($)}`;
    if (J) {
      (0, wA.checkStrictMode)(Q, `default is ignored for: ${H}`);
      return;
    }
    let B = R4._`${H} === undefined`;
    if (W.useDefaults === "empty") B = R4._`${B} || ${H} === null || ${H} === ""`;
    Y.if(B, R4._`${H} = ${(0, R4.stringify)(X)}`);
  }
});
var e0 = P((FH) => {
  Object.defineProperty(FH, "__esModule", { value: true });
  FH.validateUnion = FH.validateArray = FH.usePattern = FH.callValidateCode = FH.schemaProperties = FH.allSchemaProperties = FH.noPropertyInData = FH.propertyInData = FH.isOwnProperty = FH.hasPropFunc = FH.reportMissingProp = FH.checkMissingProp = FH.checkReportMissingProp = void 0;
  var H0 = c(), H7 = a(), Y6 = k1(), AA = a();
  function jA(Q, $) {
    let { gen: X, data: Y, it: J } = Q;
    X.if(z7(X, Y, $, J.opts.ownProperties), () => {
      Q.setParams({ missingProperty: H0._`${$}` }, true), Q.error();
    });
  }
  FH.checkReportMissingProp = jA;
  function RA({ gen: Q, data: $, it: { opts: X } }, Y, J) {
    return (0, H0.or)(...Y.map((G) => (0, H0.and)(z7(Q, $, G, X.ownProperties), H0._`${J} = ${G}`)));
  }
  FH.checkMissingProp = RA;
  function IA(Q, $) {
    Q.setParams({ missingProperty: $ }, true), Q.error();
  }
  FH.reportMissingProp = IA;
  function UH(Q) {
    return Q.scopeValue("func", { ref: Object.prototype.hasOwnProperty, code: H0._`Object.prototype.hasOwnProperty` });
  }
  FH.hasPropFunc = UH;
  function B7(Q, $, X) {
    return H0._`${UH(Q)}.call(${$}, ${X})`;
  }
  FH.isOwnProperty = B7;
  function EA(Q, $, X, Y) {
    let J = H0._`${$}${(0, H0.getProperty)(X)} !== undefined`;
    return Y ? H0._`${J} && ${B7(Q, $, X)}` : J;
  }
  FH.propertyInData = EA;
  function z7(Q, $, X, Y) {
    let J = H0._`${$}${(0, H0.getProperty)(X)} === undefined`;
    return Y ? (0, H0.or)(J, (0, H0.not)(B7(Q, $, X))) : J;
  }
  FH.noPropertyInData = z7;
  function LH(Q) {
    return Q ? Object.keys(Q).filter(($) => $ !== "__proto__") : [];
  }
  FH.allSchemaProperties = LH;
  function PA(Q, $) {
    return LH($).filter((X) => !(0, H7.alwaysValidSchema)(Q, $[X]));
  }
  FH.schemaProperties = PA;
  function bA({ schemaCode: Q, data: $, it: { gen: X, topSchemaRef: Y, schemaPath: J, errorPath: G }, it: W }, H, B, z) {
    let K = z ? H0._`${Q}, ${$}, ${Y}${J}` : $, U = [[Y6.default.instancePath, (0, H0.strConcat)(Y6.default.instancePath, G)], [Y6.default.parentData, W.parentData], [Y6.default.parentDataProperty, W.parentDataProperty], [Y6.default.rootData, Y6.default.rootData]];
    if (W.opts.dynamicRef) U.push([Y6.default.dynamicAnchors, Y6.default.dynamicAnchors]);
    let q = H0._`${K}, ${X.object(...U)}`;
    return B !== H0.nil ? H0._`${H}.call(${B}, ${q})` : H0._`${H}(${q})`;
  }
  FH.callValidateCode = bA;
  var ZA = H0._`new RegExp`;
  function CA({ gen: Q, it: { opts: $ } }, X) {
    let Y = $.unicodeRegExp ? "u" : "", { regExp: J } = $.code, G = J(X, Y);
    return Q.scopeValue("pattern", { key: G.toString(), ref: G, code: H0._`${J.code === "new RegExp" ? ZA : (0, AA.useFunc)(Q, J)}(${X}, ${Y})` });
  }
  FH.usePattern = CA;
  function SA(Q) {
    let { gen: $, data: X, keyword: Y, it: J } = Q, G = $.name("valid");
    if (J.allErrors) {
      let H = $.let("valid", true);
      return W(() => $.assign(H, false)), H;
    }
    return $.var(G, true), W(() => $.break()), G;
    function W(H) {
      let B = $.const("len", H0._`${X}.length`);
      $.forRange("i", 0, B, (z) => {
        Q.subschema({ keyword: Y, dataProp: z, dataPropType: H7.Type.Num }, G), $.if((0, H0.not)(G), H);
      });
    }
  }
  FH.validateArray = SA;
  function _A(Q) {
    let { gen: $, schema: X, keyword: Y, it: J } = Q;
    if (!Array.isArray(X)) throw Error("ajv implementation error");
    if (X.some((B) => (0, H7.alwaysValidSchema)(J, B)) && !J.opts.unevaluated) return;
    let W = $.let("valid", false), H = $.name("_valid");
    $.block(() => X.forEach((B, z) => {
      let K = Q.subschema({ keyword: Y, schemaProp: z, compositeRule: true }, H);
      if ($.assign(W, H0._`${W} || ${H}`), !Q.mergeValidEvaluated(K, H)) $.if((0, H0.not)(W));
    })), Q.result(W, () => Q.reset(), () => Q.error(true));
  }
  FH.validateUnion = _A;
});
var AH = P((wH) => {
  Object.defineProperty(wH, "__esModule", { value: true });
  wH.validateKeywordUsage = wH.validSchemaType = wH.funcKeywordCode = wH.macroKeywordCode = void 0;
  var v0 = c(), C6 = k1(), pA = e0(), dA = j9();
  function iA(Q, $) {
    let { gen: X, keyword: Y, schema: J, parentSchema: G, it: W } = Q, H = $.macro.call(W.self, J, G, W), B = DH(X, Y, H);
    if (W.opts.validateSchema !== false) W.self.validateSchema(H, true);
    let z = X.name("valid");
    Q.subschema({ schema: H, schemaPath: v0.nil, errSchemaPath: `${W.errSchemaPath}/${Y}`, topSchemaRef: B, compositeRule: true }, z), Q.pass(z, () => Q.error(true));
  }
  wH.macroKeywordCode = iA;
  function nA(Q, $) {
    var X;
    let { gen: Y, keyword: J, schema: G, parentSchema: W, $data: H, it: B } = Q;
    rA(B, $);
    let z = !H && $.compile ? $.compile.call(B.self, G, W, B) : $.validate, K = DH(Y, J, z), U = Y.let("valid");
    Q.block$data(U, q), Q.ok((X = $.valid) !== null && X !== void 0 ? X : U);
    function q() {
      if ($.errors === false) {
        if (F(), $.modifying) OH(Q);
        w(() => Q.error());
      } else {
        let D = $.async ? V() : L();
        if ($.modifying) OH(Q);
        w(() => oA(Q, D));
      }
    }
    function V() {
      let D = Y.let("ruleErrs", null);
      return Y.try(() => F(v0._`await `), (M) => Y.assign(U, false).if(v0._`${M} instanceof ${B.ValidationError}`, () => Y.assign(D, v0._`${M}.errors`), () => Y.throw(M))), D;
    }
    function L() {
      let D = v0._`${K}.errors`;
      return Y.assign(D, null), F(v0.nil), D;
    }
    function F(D = $.async ? v0._`await ` : v0.nil) {
      let M = B.opts.passContext ? C6.default.this : C6.default.self, R = !("compile" in $ && !H || $.schema === false);
      Y.assign(U, v0._`${D}${(0, pA.callValidateCode)(Q, K, M, R)}`, $.modifying);
    }
    function w(D) {
      var M;
      Y.if((0, v0.not)((M = $.valid) !== null && M !== void 0 ? M : U), D);
    }
  }
  wH.funcKeywordCode = nA;
  function OH(Q) {
    let { gen: $, data: X, it: Y } = Q;
    $.if(Y.parentData, () => $.assign(X, v0._`${Y.parentData}[${Y.parentDataProperty}]`));
  }
  function oA(Q, $) {
    let { gen: X } = Q;
    X.if(v0._`Array.isArray(${$})`, () => {
      X.assign(C6.default.vErrors, v0._`${C6.default.vErrors} === null ? ${$} : ${C6.default.vErrors}.concat(${$})`).assign(C6.default.errors, v0._`${C6.default.vErrors}.length`), (0, dA.extendErrors)(Q);
    }, () => Q.error());
  }
  function rA({ schemaEnv: Q }, $) {
    if ($.async && !Q.$async) throw Error("async keyword in sync schema");
  }
  function DH(Q, $, X) {
    if (X === void 0) throw Error(`keyword "${$}" failed to compile`);
    return Q.scopeValue("keyword", typeof X == "function" ? { ref: X } : { ref: X, code: (0, v0.stringify)(X) });
  }
  function tA(Q, $, X = false) {
    return !$.length || $.some((Y) => Y === "array" ? Array.isArray(Q) : Y === "object" ? Q && typeof Q == "object" && !Array.isArray(Q) : typeof Q == Y || X && typeof Q > "u");
  }
  wH.validSchemaType = tA;
  function aA({ schema: Q, opts: $, self: X, errSchemaPath: Y }, J, G) {
    if (Array.isArray(J.keyword) ? !J.keyword.includes(G) : J.keyword !== G) throw Error("ajv implementation error");
    let W = J.dependencies;
    if (W === null || W === void 0 ? void 0 : W.some((H) => !Object.prototype.hasOwnProperty.call(Q, H))) throw Error(`parent schema must have dependencies of ${G}: ${W.join(",")}`);
    if (J.validateSchema) {
      if (!J.validateSchema(Q[G])) {
        let B = `keyword "${G}" value is invalid at path "${Y}": ` + X.errorsText(J.validateSchema.errors);
        if ($.validateSchema === "log") X.logger.error(B);
        else throw Error(B);
      }
    }
  }
  wH.validateKeywordUsage = aA;
});
var EH = P((RH) => {
  Object.defineProperty(RH, "__esModule", { value: true });
  RH.extendSubschemaMode = RH.extendSubschemaData = RH.getSubschema = void 0;
  var R1 = c(), jH = a();
  function $j(Q, { keyword: $, schemaProp: X, schema: Y, schemaPath: J, errSchemaPath: G, topSchemaRef: W }) {
    if ($ !== void 0 && Y !== void 0) throw Error('both "keyword" and "schema" passed, only one allowed');
    if ($ !== void 0) {
      let H = Q.schema[$];
      return X === void 0 ? { schema: H, schemaPath: R1._`${Q.schemaPath}${(0, R1.getProperty)($)}`, errSchemaPath: `${Q.errSchemaPath}/${$}` } : { schema: H[X], schemaPath: R1._`${Q.schemaPath}${(0, R1.getProperty)($)}${(0, R1.getProperty)(X)}`, errSchemaPath: `${Q.errSchemaPath}/${$}/${(0, jH.escapeFragment)(X)}` };
    }
    if (Y !== void 0) {
      if (J === void 0 || G === void 0 || W === void 0) throw Error('"schemaPath", "errSchemaPath" and "topSchemaRef" are required with "schema"');
      return { schema: Y, schemaPath: J, topSchemaRef: W, errSchemaPath: G };
    }
    throw Error('either "keyword" or "schema" must be passed');
  }
  RH.getSubschema = $j;
  function Xj(Q, $, { dataProp: X, dataPropType: Y, data: J, dataTypes: G, propertyName: W }) {
    if (J !== void 0 && X !== void 0) throw Error('both "data" and "dataProp" passed, only one allowed');
    let { gen: H } = $;
    if (X !== void 0) {
      let { errorPath: z, dataPathArr: K, opts: U } = $, q = H.let("data", R1._`${$.data}${(0, R1.getProperty)(X)}`, true);
      B(q), Q.errorPath = R1.str`${z}${(0, jH.getErrorPath)(X, Y, U.jsPropertySyntax)}`, Q.parentDataProperty = R1._`${X}`, Q.dataPathArr = [...K, Q.parentDataProperty];
    }
    if (J !== void 0) {
      let z = J instanceof R1.Name ? J : H.let("data", J, true);
      if (B(z), W !== void 0) Q.propertyName = W;
    }
    if (G) Q.dataTypes = G;
    function B(z) {
      Q.data = z, Q.dataLevel = $.dataLevel + 1, Q.dataTypes = [], $.definedProperties = /* @__PURE__ */ new Set(), Q.parentData = $.data, Q.dataNames = [...$.dataNames, z];
    }
  }
  RH.extendSubschemaData = Xj;
  function Yj(Q, { jtdDiscriminator: $, jtdMetadata: X, compositeRule: Y, createErrors: J, allErrors: G }) {
    if (Y !== void 0) Q.compositeRule = Y;
    if (J !== void 0) Q.createErrors = J;
    if (G !== void 0) Q.allErrors = G;
    Q.jtdDiscriminator = $, Q.jtdMetadata = X;
  }
  RH.extendSubschemaMode = Yj;
});
var K7 = P((Ly, PH) => {
  PH.exports = function Q($, X) {
    if ($ === X) return true;
    if ($ && X && typeof $ == "object" && typeof X == "object") {
      if ($.constructor !== X.constructor) return false;
      var Y, J, G;
      if (Array.isArray($)) {
        if (Y = $.length, Y != X.length) return false;
        for (J = Y; J-- !== 0; ) if (!Q($[J], X[J])) return false;
        return true;
      }
      if ($.constructor === RegExp) return $.source === X.source && $.flags === X.flags;
      if ($.valueOf !== Object.prototype.valueOf) return $.valueOf() === X.valueOf();
      if ($.toString !== Object.prototype.toString) return $.toString() === X.toString();
      if (G = Object.keys($), Y = G.length, Y !== Object.keys(X).length) return false;
      for (J = Y; J-- !== 0; ) if (!Object.prototype.hasOwnProperty.call(X, G[J])) return false;
      for (J = Y; J-- !== 0; ) {
        var W = G[J];
        if (!Q($[W], X[W])) return false;
      }
      return true;
    }
    return $ !== $ && X !== X;
  };
});
var ZH = P((Fy, bH) => {
  var J6 = bH.exports = function(Q, $, X) {
    if (typeof $ == "function") X = $, $ = {};
    X = $.cb || X;
    var Y = typeof X == "function" ? X : X.pre || function() {
    }, J = X.post || function() {
    };
    qQ($, Y, J, Q, "", Q);
  };
  J6.keywords = { additionalItems: true, items: true, contains: true, additionalProperties: true, propertyNames: true, not: true, if: true, then: true, else: true };
  J6.arrayKeywords = { items: true, allOf: true, anyOf: true, oneOf: true };
  J6.propsKeywords = { $defs: true, definitions: true, properties: true, patternProperties: true, dependencies: true };
  J6.skipKeywords = { default: true, enum: true, const: true, required: true, maximum: true, minimum: true, exclusiveMaximum: true, exclusiveMinimum: true, multipleOf: true, maxLength: true, minLength: true, pattern: true, format: true, maxItems: true, minItems: true, uniqueItems: true, maxProperties: true, minProperties: true };
  function qQ(Q, $, X, Y, J, G, W, H, B, z) {
    if (Y && typeof Y == "object" && !Array.isArray(Y)) {
      $(Y, J, G, W, H, B, z);
      for (var K in Y) {
        var U = Y[K];
        if (Array.isArray(U)) {
          if (K in J6.arrayKeywords) for (var q = 0; q < U.length; q++) qQ(Q, $, X, U[q], J + "/" + K + "/" + q, G, J, K, Y, q);
        } else if (K in J6.propsKeywords) {
          if (U && typeof U == "object") for (var V in U) qQ(Q, $, X, U[V], J + "/" + K + "/" + Wj(V), G, J, K, Y, V);
        } else if (K in J6.keywords || Q.allKeys && !(K in J6.skipKeywords)) qQ(Q, $, X, U, J + "/" + K, G, J, K, Y);
      }
      X(Y, J, G, W, H, B, z);
    }
  }
  function Wj(Q) {
    return Q.replace(/~/g, "~0").replace(/\//g, "~1");
  }
});
var I9 = P((kH) => {
  Object.defineProperty(kH, "__esModule", { value: true });
  kH.getSchemaRefs = kH.resolveUrl = kH.normalizeId = kH._getFullPath = kH.getFullPath = kH.inlineRef = void 0;
  var Hj = a(), Bj = K7(), zj = ZH(), Kj = /* @__PURE__ */ new Set(["type", "format", "pattern", "maxLength", "minLength", "maxProperties", "minProperties", "maxItems", "minItems", "maximum", "minimum", "uniqueItems", "multipleOf", "required", "enum", "const"]);
  function Vj(Q, $ = true) {
    if (typeof Q == "boolean") return true;
    if ($ === true) return !V7(Q);
    if (!$) return false;
    return CH(Q) <= $;
  }
  kH.inlineRef = Vj;
  var qj = /* @__PURE__ */ new Set(["$ref", "$recursiveRef", "$recursiveAnchor", "$dynamicRef", "$dynamicAnchor"]);
  function V7(Q) {
    for (let $ in Q) {
      if (qj.has($)) return true;
      let X = Q[$];
      if (Array.isArray(X) && X.some(V7)) return true;
      if (typeof X == "object" && V7(X)) return true;
    }
    return false;
  }
  function CH(Q) {
    let $ = 0;
    for (let X in Q) {
      if (X === "$ref") return 1 / 0;
      if ($++, Kj.has(X)) continue;
      if (typeof Q[X] == "object") (0, Hj.eachItem)(Q[X], (Y) => $ += CH(Y));
      if ($ === 1 / 0) return 1 / 0;
    }
    return $;
  }
  function SH(Q, $ = "", X) {
    if (X !== false) $ = I4($);
    let Y = Q.parse($);
    return _H(Q, Y);
  }
  kH.getFullPath = SH;
  function _H(Q, $) {
    return Q.serialize($).split("#")[0] + "#";
  }
  kH._getFullPath = _H;
  var Uj = /#\/?$/;
  function I4(Q) {
    return Q ? Q.replace(Uj, "") : "";
  }
  kH.normalizeId = I4;
  function Lj(Q, $, X) {
    return X = I4(X), Q.resolve($, X);
  }
  kH.resolveUrl = Lj;
  var Fj = /^[a-z_][-a-z0-9._]*$/i;
  function Nj(Q, $) {
    if (typeof Q == "boolean") return {};
    let { schemaId: X, uriResolver: Y } = this.opts, J = I4(Q[X] || $), G = { "": J }, W = SH(Y, J, false), H = {}, B = /* @__PURE__ */ new Set();
    return zj(Q, { allKeys: true }, (U, q, V, L) => {
      if (L === void 0) return;
      let F = W + q, w = G[L];
      if (typeof U[X] == "string") w = D.call(this, U[X]);
      M.call(this, U.$anchor), M.call(this, U.$dynamicAnchor), G[q] = w;
      function D(R) {
        let Z = this.opts.uriResolver.resolve;
        if (R = I4(w ? Z(w, R) : R), B.has(R)) throw K(R);
        B.add(R);
        let v = this.refs[R];
        if (typeof v == "string") v = this.refs[v];
        if (typeof v == "object") z(U, v.schema, R);
        else if (R !== I4(F)) if (R[0] === "#") z(U, H[R], R), H[R] = U;
        else this.refs[R] = F;
        return R;
      }
      function M(R) {
        if (typeof R == "string") {
          if (!Fj.test(R)) throw Error(`invalid anchor "${R}"`);
          D.call(this, `#${R}`);
        }
      }
    }), H;
    function z(U, q, V) {
      if (q !== void 0 && !Bj(U, q)) throw K(V);
    }
    function K(U) {
      return Error(`reference "${U}" resolves to more than one schema`);
    }
  }
  kH.getSchemaRefs = Nj;
});
var b9 = P((oH) => {
  Object.defineProperty(oH, "__esModule", { value: true });
  oH.getData = oH.KeywordCxt = oH.validateFunctionCode = void 0;
  var hH = a3(), TH = R9(), U7 = Y7(), UQ = R9(), jj = qH(), P9 = AH(), q7 = EH(), k = c(), f = k1(), Rj = I9(), v1 = a(), E9 = j9();
  function Ij(Q) {
    if (mH(Q)) {
      if (lH(Q), uH(Q)) {
        bj(Q);
        return;
      }
    }
    fH(Q, () => (0, hH.topBoolOrEmptySchema)(Q));
  }
  oH.validateFunctionCode = Ij;
  function fH({ gen: Q, validateName: $, schema: X, schemaEnv: Y, opts: J }, G) {
    if (J.code.es5) Q.func($, k._`${f.default.data}, ${f.default.valCxt}`, Y.$async, () => {
      Q.code(k._`"use strict"; ${xH(X, J)}`), Pj(Q, J), Q.code(G);
    });
    else Q.func($, k._`${f.default.data}, ${Ej(J)}`, Y.$async, () => Q.code(xH(X, J)).code(G));
  }
  function Ej(Q) {
    return k._`{${f.default.instancePath}="", ${f.default.parentData}, ${f.default.parentDataProperty}, ${f.default.rootData}=${f.default.data}${Q.dynamicRef ? k._`, ${f.default.dynamicAnchors}={}` : k.nil}}={}`;
  }
  function Pj(Q, $) {
    Q.if(f.default.valCxt, () => {
      if (Q.var(f.default.instancePath, k._`${f.default.valCxt}.${f.default.instancePath}`), Q.var(f.default.parentData, k._`${f.default.valCxt}.${f.default.parentData}`), Q.var(f.default.parentDataProperty, k._`${f.default.valCxt}.${f.default.parentDataProperty}`), Q.var(f.default.rootData, k._`${f.default.valCxt}.${f.default.rootData}`), $.dynamicRef) Q.var(f.default.dynamicAnchors, k._`${f.default.valCxt}.${f.default.dynamicAnchors}`);
    }, () => {
      if (Q.var(f.default.instancePath, k._`""`), Q.var(f.default.parentData, k._`undefined`), Q.var(f.default.parentDataProperty, k._`undefined`), Q.var(f.default.rootData, f.default.data), $.dynamicRef) Q.var(f.default.dynamicAnchors, k._`{}`);
    });
  }
  function bj(Q) {
    let { schema: $, opts: X, gen: Y } = Q;
    fH(Q, () => {
      if (X.$comment && $.$comment) pH(Q);
      if (kj(Q), Y.let(f.default.vErrors, null), Y.let(f.default.errors, 0), X.unevaluated) Zj(Q);
      cH(Q), xj(Q);
    });
    return;
  }
  function Zj(Q) {
    let { gen: $, validateName: X } = Q;
    Q.evaluated = $.const("evaluated", k._`${X}.evaluated`), $.if(k._`${Q.evaluated}.dynamicProps`, () => $.assign(k._`${Q.evaluated}.props`, k._`undefined`)), $.if(k._`${Q.evaluated}.dynamicItems`, () => $.assign(k._`${Q.evaluated}.items`, k._`undefined`));
  }
  function xH(Q, $) {
    let X = typeof Q == "object" && Q[$.schemaId];
    return X && ($.code.source || $.code.process) ? k._`/*# sourceURL=${X} */` : k.nil;
  }
  function Cj(Q, $) {
    if (mH(Q)) {
      if (lH(Q), uH(Q)) {
        Sj(Q, $);
        return;
      }
    }
    (0, hH.boolOrEmptySchema)(Q, $);
  }
  function uH({ schema: Q, self: $ }) {
    if (typeof Q == "boolean") return !Q;
    for (let X in Q) if ($.RULES.all[X]) return true;
    return false;
  }
  function mH(Q) {
    return typeof Q.schema != "boolean";
  }
  function Sj(Q, $) {
    let { schema: X, gen: Y, opts: J } = Q;
    if (J.$comment && X.$comment) pH(Q);
    vj(Q), Tj(Q);
    let G = Y.const("_errs", f.default.errors);
    cH(Q, G), Y.var($, k._`${G} === ${f.default.errors}`);
  }
  function lH(Q) {
    (0, v1.checkUnknownRules)(Q), _j(Q);
  }
  function cH(Q, $) {
    if (Q.opts.jtd) return yH(Q, [], false, $);
    let X = (0, TH.getSchemaTypes)(Q.schema), Y = (0, TH.coerceAndCheckDataType)(Q, X);
    yH(Q, X, !Y, $);
  }
  function _j(Q) {
    let { schema: $, errSchemaPath: X, opts: Y, self: J } = Q;
    if ($.$ref && Y.ignoreKeywordsWithRef && (0, v1.schemaHasRulesButRef)($, J.RULES)) J.logger.warn(`$ref: keywords ignored in schema at path "${X}"`);
  }
  function kj(Q) {
    let { schema: $, opts: X } = Q;
    if ($.default !== void 0 && X.useDefaults && X.strictSchema) (0, v1.checkStrictMode)(Q, "default is ignored in the schema root");
  }
  function vj(Q) {
    let $ = Q.schema[Q.opts.schemaId];
    if ($) Q.baseId = (0, Rj.resolveUrl)(Q.opts.uriResolver, Q.baseId, $);
  }
  function Tj(Q) {
    if (Q.schema.$async && !Q.schemaEnv.$async) throw Error("async schema in sync schema");
  }
  function pH({ gen: Q, schemaEnv: $, schema: X, errSchemaPath: Y, opts: J }) {
    let G = X.$comment;
    if (J.$comment === true) Q.code(k._`${f.default.self}.logger.log(${G})`);
    else if (typeof J.$comment == "function") {
      let W = k.str`${Y}/$comment`, H = Q.scopeValue("root", { ref: $.root });
      Q.code(k._`${f.default.self}.opts.$comment(${G}, ${W}, ${H}.schema)`);
    }
  }
  function xj(Q) {
    let { gen: $, schemaEnv: X, validateName: Y, ValidationError: J, opts: G } = Q;
    if (X.$async) $.if(k._`${f.default.errors} === 0`, () => $.return(f.default.data), () => $.throw(k._`new ${J}(${f.default.vErrors})`));
    else {
      if ($.assign(k._`${Y}.errors`, f.default.vErrors), G.unevaluated) yj(Q);
      $.return(k._`${f.default.errors} === 0`);
    }
  }
  function yj({ gen: Q, evaluated: $, props: X, items: Y }) {
    if (X instanceof k.Name) Q.assign(k._`${$}.props`, X);
    if (Y instanceof k.Name) Q.assign(k._`${$}.items`, Y);
  }
  function yH(Q, $, X, Y) {
    let { gen: J, schema: G, data: W, allErrors: H, opts: B, self: z } = Q, { RULES: K } = z;
    if (G.$ref && (B.ignoreKeywordsWithRef || !(0, v1.schemaHasRulesButRef)(G, K))) {
      J.block(() => iH(Q, "$ref", K.all.$ref.definition));
      return;
    }
    if (!B.jtd) gj(Q, $);
    J.block(() => {
      for (let q of K.rules) U(q);
      U(K.post);
    });
    function U(q) {
      if (!(0, U7.shouldUseGroup)(G, q)) return;
      if (q.type) {
        if (J.if((0, UQ.checkDataType)(q.type, W, B.strictNumbers)), gH(Q, q), $.length === 1 && $[0] === q.type && X) J.else(), (0, UQ.reportTypeError)(Q);
        J.endIf();
      } else gH(Q, q);
      if (!H) J.if(k._`${f.default.errors} === ${Y || 0}`);
    }
  }
  function gH(Q, $) {
    let { gen: X, schema: Y, opts: { useDefaults: J } } = Q;
    if (J) (0, jj.assignDefaults)(Q, $.type);
    X.block(() => {
      for (let G of $.rules) if ((0, U7.shouldUseRule)(Y, G)) iH(Q, G.keyword, G.definition, $.type);
    });
  }
  function gj(Q, $) {
    if (Q.schemaEnv.meta || !Q.opts.strictTypes) return;
    if (hj(Q, $), !Q.opts.allowUnionTypes) fj(Q, $);
    uj(Q, Q.dataTypes);
  }
  function hj(Q, $) {
    if (!$.length) return;
    if (!Q.dataTypes.length) {
      Q.dataTypes = $;
      return;
    }
    $.forEach((X) => {
      if (!dH(Q.dataTypes, X)) L7(Q, `type "${X}" not allowed by context "${Q.dataTypes.join(",")}"`);
    }), lj(Q, $);
  }
  function fj(Q, $) {
    if ($.length > 1 && !($.length === 2 && $.includes("null"))) L7(Q, "use allowUnionTypes to allow union type keyword");
  }
  function uj(Q, $) {
    let X = Q.self.RULES.all;
    for (let Y in X) {
      let J = X[Y];
      if (typeof J == "object" && (0, U7.shouldUseRule)(Q.schema, J)) {
        let { type: G } = J.definition;
        if (G.length && !G.some((W) => mj($, W))) L7(Q, `missing type "${G.join(",")}" for keyword "${Y}"`);
      }
    }
  }
  function mj(Q, $) {
    return Q.includes($) || $ === "number" && Q.includes("integer");
  }
  function dH(Q, $) {
    return Q.includes($) || $ === "integer" && Q.includes("number");
  }
  function lj(Q, $) {
    let X = [];
    for (let Y of Q.dataTypes) if (dH($, Y)) X.push(Y);
    else if ($.includes("integer") && Y === "number") X.push("integer");
    Q.dataTypes = X;
  }
  function L7(Q, $) {
    let X = Q.schemaEnv.baseId + Q.errSchemaPath;
    $ += ` at "${X}" (strictTypes)`, (0, v1.checkStrictMode)(Q, $, Q.opts.strictTypes);
  }
  class F7 {
    constructor(Q, $, X) {
      if ((0, P9.validateKeywordUsage)(Q, $, X), this.gen = Q.gen, this.allErrors = Q.allErrors, this.keyword = X, this.data = Q.data, this.schema = Q.schema[X], this.$data = $.$data && Q.opts.$data && this.schema && this.schema.$data, this.schemaValue = (0, v1.schemaRefOrVal)(Q, this.schema, X, this.$data), this.schemaType = $.schemaType, this.parentSchema = Q.schema, this.params = {}, this.it = Q, this.def = $, this.$data) this.schemaCode = Q.gen.const("vSchema", nH(this.$data, Q));
      else if (this.schemaCode = this.schemaValue, !(0, P9.validSchemaType)(this.schema, $.schemaType, $.allowUndefined)) throw Error(`${X} value must be ${JSON.stringify($.schemaType)}`);
      if ("code" in $ ? $.trackErrors : $.errors !== false) this.errsCount = Q.gen.const("_errs", f.default.errors);
    }
    result(Q, $, X) {
      this.failResult((0, k.not)(Q), $, X);
    }
    failResult(Q, $, X) {
      if (this.gen.if(Q), X) X();
      else this.error();
      if ($) {
        if (this.gen.else(), $(), this.allErrors) this.gen.endIf();
      } else if (this.allErrors) this.gen.endIf();
      else this.gen.else();
    }
    pass(Q, $) {
      this.failResult((0, k.not)(Q), void 0, $);
    }
    fail(Q) {
      if (Q === void 0) {
        if (this.error(), !this.allErrors) this.gen.if(false);
        return;
      }
      if (this.gen.if(Q), this.error(), this.allErrors) this.gen.endIf();
      else this.gen.else();
    }
    fail$data(Q) {
      if (!this.$data) return this.fail(Q);
      let { schemaCode: $ } = this;
      this.fail(k._`${$} !== undefined && (${(0, k.or)(this.invalid$data(), Q)})`);
    }
    error(Q, $, X) {
      if ($) {
        this.setParams($), this._error(Q, X), this.setParams({});
        return;
      }
      this._error(Q, X);
    }
    _error(Q, $) {
      (Q ? E9.reportExtraError : E9.reportError)(this, this.def.error, $);
    }
    $dataError() {
      (0, E9.reportError)(this, this.def.$dataError || E9.keyword$DataError);
    }
    reset() {
      if (this.errsCount === void 0) throw Error('add "trackErrors" to keyword definition');
      (0, E9.resetErrorsCount)(this.gen, this.errsCount);
    }
    ok(Q) {
      if (!this.allErrors) this.gen.if(Q);
    }
    setParams(Q, $) {
      if ($) Object.assign(this.params, Q);
      else this.params = Q;
    }
    block$data(Q, $, X = k.nil) {
      this.gen.block(() => {
        this.check$data(Q, X), $();
      });
    }
    check$data(Q = k.nil, $ = k.nil) {
      if (!this.$data) return;
      let { gen: X, schemaCode: Y, schemaType: J, def: G } = this;
      if (X.if((0, k.or)(k._`${Y} === undefined`, $)), Q !== k.nil) X.assign(Q, true);
      if (J.length || G.validateSchema) {
        if (X.elseIf(this.invalid$data()), this.$dataError(), Q !== k.nil) X.assign(Q, false);
      }
      X.else();
    }
    invalid$data() {
      let { gen: Q, schemaCode: $, schemaType: X, def: Y, it: J } = this;
      return (0, k.or)(G(), W());
      function G() {
        if (X.length) {
          if (!($ instanceof k.Name)) throw Error("ajv implementation error");
          let H = Array.isArray(X) ? X : [X];
          return k._`${(0, UQ.checkDataTypes)(H, $, J.opts.strictNumbers, UQ.DataType.Wrong)}`;
        }
        return k.nil;
      }
      function W() {
        if (Y.validateSchema) {
          let H = Q.scopeValue("validate$data", { ref: Y.validateSchema });
          return k._`!${H}(${$})`;
        }
        return k.nil;
      }
    }
    subschema(Q, $) {
      let X = (0, q7.getSubschema)(this.it, Q);
      (0, q7.extendSubschemaData)(X, this.it, Q), (0, q7.extendSubschemaMode)(X, Q);
      let Y = { ...this.it, ...X, items: void 0, props: void 0 };
      return Cj(Y, $), Y;
    }
    mergeEvaluated(Q, $) {
      let { it: X, gen: Y } = this;
      if (!X.opts.unevaluated) return;
      if (X.props !== true && Q.props !== void 0) X.props = v1.mergeEvaluated.props(Y, Q.props, X.props, $);
      if (X.items !== true && Q.items !== void 0) X.items = v1.mergeEvaluated.items(Y, Q.items, X.items, $);
    }
    mergeValidEvaluated(Q, $) {
      let { it: X, gen: Y } = this;
      if (X.opts.unevaluated && (X.props !== true || X.items !== true)) return Y.if($, () => this.mergeEvaluated(Q, k.Name)), true;
    }
  }
  oH.KeywordCxt = F7;
  function iH(Q, $, X, Y) {
    let J = new F7(Q, X, $);
    if ("code" in X) X.code(J, Y);
    else if (J.$data && X.validate) (0, P9.funcKeywordCode)(J, X);
    else if ("macro" in X) (0, P9.macroKeywordCode)(J, X);
    else if (X.compile || X.validate) (0, P9.funcKeywordCode)(J, X);
  }
  var cj = /^\/(?:[^~]|~0|~1)*$/, pj = /^([0-9]+)(#|\/(?:[^~]|~0|~1)*)?$/;
  function nH(Q, { dataLevel: $, dataNames: X, dataPathArr: Y }) {
    let J, G;
    if (Q === "") return f.default.rootData;
    if (Q[0] === "/") {
      if (!cj.test(Q)) throw Error(`Invalid JSON-pointer: ${Q}`);
      J = Q, G = f.default.rootData;
    } else {
      let z = pj.exec(Q);
      if (!z) throw Error(`Invalid JSON-pointer: ${Q}`);
      let K = +z[1];
      if (J = z[2], J === "#") {
        if (K >= $) throw Error(B("property/index", K));
        return Y[$ - K];
      }
      if (K > $) throw Error(B("data", K));
      if (G = X[$ - K], !J) return G;
    }
    let W = G, H = J.split("/");
    for (let z of H) if (z) G = k._`${G}${(0, k.getProperty)((0, v1.unescapeJsonPointer)(z))}`, W = k._`${W} && ${G}`;
    return W;
    function B(z, K) {
      return `Cannot access ${z} ${K} levels up, current level is ${$}`;
    }
  }
  oH.getData = nH;
});
var LQ = P((aH) => {
  Object.defineProperty(aH, "__esModule", { value: true });
  class tH extends Error {
    constructor(Q) {
      super("validation failed");
      this.errors = Q, this.ajv = this.validation = true;
    }
  }
  aH.default = tH;
});
var Z9 = P((eH) => {
  Object.defineProperty(eH, "__esModule", { value: true });
  var N7 = I9();
  class sH extends Error {
    constructor(Q, $, X, Y) {
      super(Y || `can't resolve reference ${X} from id ${$}`);
      this.missingRef = (0, N7.resolveUrl)(Q, $, X), this.missingSchema = (0, N7.normalizeId)((0, N7.getFullPath)(Q, this.missingRef));
    }
  }
  eH.default = sH;
});
var NQ = P((XB) => {
  Object.defineProperty(XB, "__esModule", { value: true });
  XB.resolveSchema = XB.getCompilingSchema = XB.resolveRef = XB.compileSchema = XB.SchemaEnv = void 0;
  var V1 = c(), rj = LQ(), S6 = k1(), q1 = I9(), QB = a(), tj = b9();
  class C9 {
    constructor(Q) {
      var $;
      this.refs = {}, this.dynamicAnchors = {};
      let X;
      if (typeof Q.schema == "object") X = Q.schema;
      this.schema = Q.schema, this.schemaId = Q.schemaId, this.root = Q.root || this, this.baseId = ($ = Q.baseId) !== null && $ !== void 0 ? $ : (0, q1.normalizeId)(X === null || X === void 0 ? void 0 : X[Q.schemaId || "$id"]), this.schemaPath = Q.schemaPath, this.localRefs = Q.localRefs, this.meta = Q.meta, this.$async = X === null || X === void 0 ? void 0 : X.$async, this.refs = {};
    }
  }
  XB.SchemaEnv = C9;
  function D7(Q) {
    let $ = $B.call(this, Q);
    if ($) return $;
    let X = (0, q1.getFullPath)(this.opts.uriResolver, Q.root.baseId), { es5: Y, lines: J } = this.opts.code, { ownProperties: G } = this.opts, W = new V1.CodeGen(this.scope, { es5: Y, lines: J, ownProperties: G }), H;
    if (Q.$async) H = W.scopeValue("Error", { ref: rj.default, code: V1._`require("ajv/dist/runtime/validation_error").default` });
    let B = W.scopeName("validate");
    Q.validateName = B;
    let z = { gen: W, allErrors: this.opts.allErrors, data: S6.default.data, parentData: S6.default.parentData, parentDataProperty: S6.default.parentDataProperty, dataNames: [S6.default.data], dataPathArr: [V1.nil], dataLevel: 0, dataTypes: [], definedProperties: /* @__PURE__ */ new Set(), topSchemaRef: W.scopeValue("schema", this.opts.code.source === true ? { ref: Q.schema, code: (0, V1.stringify)(Q.schema) } : { ref: Q.schema }), validateName: B, ValidationError: H, schema: Q.schema, schemaEnv: Q, rootId: X, baseId: Q.baseId || X, schemaPath: V1.nil, errSchemaPath: Q.schemaPath || (this.opts.jtd ? "" : "#"), errorPath: V1._`""`, opts: this.opts, self: this }, K;
    try {
      this._compilations.add(Q), (0, tj.validateFunctionCode)(z), W.optimize(this.opts.code.optimize);
      let U = W.toString();
      if (K = `${W.scopeRefs(S6.default.scope)}return ${U}`, this.opts.code.process) K = this.opts.code.process(K, Q);
      let V = Function(`${S6.default.self}`, `${S6.default.scope}`, K)(this, this.scope.get());
      if (this.scope.value(B, { ref: V }), V.errors = null, V.schema = Q.schema, V.schemaEnv = Q, Q.$async) V.$async = true;
      if (this.opts.code.source === true) V.source = { validateName: B, validateCode: U, scopeValues: W._values };
      if (this.opts.unevaluated) {
        let { props: L, items: F } = z;
        if (V.evaluated = { props: L instanceof V1.Name ? void 0 : L, items: F instanceof V1.Name ? void 0 : F, dynamicProps: L instanceof V1.Name, dynamicItems: F instanceof V1.Name }, V.source) V.source.evaluated = (0, V1.stringify)(V.evaluated);
      }
      return Q.validate = V, Q;
    } catch (U) {
      if (delete Q.validate, delete Q.validateName, K) this.logger.error("Error compiling schema, function code:", K);
      throw U;
    } finally {
      this._compilations.delete(Q);
    }
  }
  XB.compileSchema = D7;
  function aj(Q, $, X) {
    var Y;
    X = (0, q1.resolveUrl)(this.opts.uriResolver, $, X);
    let J = Q.refs[X];
    if (J) return J;
    let G = Q2.call(this, Q, X);
    if (G === void 0) {
      let W = (Y = Q.localRefs) === null || Y === void 0 ? void 0 : Y[X], { schemaId: H } = this.opts;
      if (W) G = new C9({ schema: W, schemaId: H, root: Q, baseId: $ });
    }
    if (G === void 0) return;
    return Q.refs[X] = sj.call(this, G);
  }
  XB.resolveRef = aj;
  function sj(Q) {
    if ((0, q1.inlineRef)(Q.schema, this.opts.inlineRefs)) return Q.schema;
    return Q.validate ? Q : D7.call(this, Q);
  }
  function $B(Q) {
    for (let $ of this._compilations) if (ej($, Q)) return $;
  }
  XB.getCompilingSchema = $B;
  function ej(Q, $) {
    return Q.schema === $.schema && Q.root === $.root && Q.baseId === $.baseId;
  }
  function Q2(Q, $) {
    let X;
    while (typeof (X = this.refs[$]) == "string") $ = X;
    return X || this.schemas[$] || FQ.call(this, Q, $);
  }
  function FQ(Q, $) {
    let X = this.opts.uriResolver.parse($), Y = (0, q1._getFullPath)(this.opts.uriResolver, X), J = (0, q1.getFullPath)(this.opts.uriResolver, Q.baseId, void 0);
    if (Object.keys(Q.schema).length > 0 && Y === J) return O7.call(this, X, Q);
    let G = (0, q1.normalizeId)(Y), W = this.refs[G] || this.schemas[G];
    if (typeof W == "string") {
      let H = FQ.call(this, Q, W);
      if (typeof (H === null || H === void 0 ? void 0 : H.schema) !== "object") return;
      return O7.call(this, X, H);
    }
    if (typeof (W === null || W === void 0 ? void 0 : W.schema) !== "object") return;
    if (!W.validate) D7.call(this, W);
    if (G === (0, q1.normalizeId)($)) {
      let { schema: H } = W, { schemaId: B } = this.opts, z = H[B];
      if (z) J = (0, q1.resolveUrl)(this.opts.uriResolver, J, z);
      return new C9({ schema: H, schemaId: B, root: Q, baseId: J });
    }
    return O7.call(this, X, W);
  }
  XB.resolveSchema = FQ;
  var $2 = /* @__PURE__ */ new Set(["properties", "patternProperties", "enum", "dependencies", "definitions"]);
  function O7(Q, { baseId: $, schema: X, root: Y }) {
    var J;
    if (((J = Q.fragment) === null || J === void 0 ? void 0 : J[0]) !== "/") return;
    for (let H of Q.fragment.slice(1).split("/")) {
      if (typeof X === "boolean") return;
      let B = X[(0, QB.unescapeFragment)(H)];
      if (B === void 0) return;
      X = B;
      let z = typeof X === "object" && X[this.opts.schemaId];
      if (!$2.has(H) && z) $ = (0, q1.resolveUrl)(this.opts.uriResolver, $, z);
    }
    let G;
    if (typeof X != "boolean" && X.$ref && !(0, QB.schemaHasRulesButRef)(X, this.RULES)) {
      let H = (0, q1.resolveUrl)(this.opts.uriResolver, $, X.$ref);
      G = FQ.call(this, Y, H);
    }
    let { schemaId: W } = this.opts;
    if (G = G || new C9({ schema: X, schemaId: W, root: Y, baseId: $ }), G.schema !== G.root.schema) return G;
    return;
  }
});
var JB = P((Ay, W2) => {
  W2.exports = { $id: "https://raw.githubusercontent.com/ajv-validator/ajv/master/lib/refs/data.json#", description: "Meta-schema for $data reference (JSON AnySchema extension proposal)", type: "object", required: ["$data"], properties: { $data: { type: "string", anyOf: [{ format: "relative-json-pointer" }, { format: "json-pointer" }] } }, additionalProperties: false };
});
var WB = P((jy, GB) => {
  var H2 = { 0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, a: 10, A: 10, b: 11, B: 11, c: 12, C: 12, d: 13, D: 13, e: 14, E: 14, f: 15, F: 15 };
  GB.exports = { HEX: H2 };
});
var LB = P((Ry, UB) => {
  var { HEX: B2 } = WB(), z2 = /^(?:(?:25[0-5]|2[0-4]\d|1\d{2}|[1-9]\d|\d)\.){3}(?:25[0-5]|2[0-4]\d|1\d{2}|[1-9]\d|\d)$/u;
  function KB(Q) {
    if (qB(Q, ".") < 3) return { host: Q, isIPV4: false };
    let $ = Q.match(z2) || [], [X] = $;
    if (X) return { host: V2(X, "."), isIPV4: true };
    else return { host: Q, isIPV4: false };
  }
  function w7(Q, $ = false) {
    let X = "", Y = true;
    for (let J of Q) {
      if (B2[J] === void 0) return;
      if (J !== "0" && Y === true) Y = false;
      if (!Y) X += J;
    }
    if ($ && X.length === 0) X = "0";
    return X;
  }
  function K2(Q) {
    let $ = 0, X = { error: false, address: "", zone: "" }, Y = [], J = [], G = false, W = false, H = false;
    function B() {
      if (J.length) {
        if (G === false) {
          let z = w7(J);
          if (z !== void 0) Y.push(z);
          else return X.error = true, false;
        }
        J.length = 0;
      }
      return true;
    }
    for (let z = 0; z < Q.length; z++) {
      let K = Q[z];
      if (K === "[" || K === "]") continue;
      if (K === ":") {
        if (W === true) H = true;
        if (!B()) break;
        if ($++, Y.push(":"), $ > 7) {
          X.error = true;
          break;
        }
        if (z - 1 >= 0 && Q[z - 1] === ":") W = true;
        continue;
      } else if (K === "%") {
        if (!B()) break;
        G = true;
      } else {
        J.push(K);
        continue;
      }
    }
    if (J.length) if (G) X.zone = J.join("");
    else if (H) Y.push(J.join(""));
    else Y.push(w7(J));
    return X.address = Y.join(""), X;
  }
  function VB(Q) {
    if (qB(Q, ":") < 2) return { host: Q, isIPV6: false };
    let $ = K2(Q);
    if (!$.error) {
      let { address: X, address: Y } = $;
      if ($.zone) X += "%" + $.zone, Y += "%25" + $.zone;
      return { host: X, escapedHost: Y, isIPV6: true };
    } else return { host: Q, isIPV6: false };
  }
  function V2(Q, $) {
    let X = "", Y = true, J = Q.length;
    for (let G = 0; G < J; G++) {
      let W = Q[G];
      if (W === "0" && Y) {
        if (G + 1 <= J && Q[G + 1] === $ || G + 1 === J) X += W, Y = false;
      } else {
        if (W === $) Y = true;
        else Y = false;
        X += W;
      }
    }
    return X;
  }
  function qB(Q, $) {
    let X = 0;
    for (let Y = 0; Y < Q.length; Y++) if (Q[Y] === $) X++;
    return X;
  }
  var HB = /^\.\.?\//u, BB = /^\/\.(?:\/|$)/u, zB = /^\/\.\.(?:\/|$)/u, q2 = /^\/?(?:.|\n)*?(?=\/|$)/u;
  function U2(Q) {
    let $ = [];
    while (Q.length) if (Q.match(HB)) Q = Q.replace(HB, "");
    else if (Q.match(BB)) Q = Q.replace(BB, "/");
    else if (Q.match(zB)) Q = Q.replace(zB, "/"), $.pop();
    else if (Q === "." || Q === "..") Q = "";
    else {
      let X = Q.match(q2);
      if (X) {
        let Y = X[0];
        Q = Q.slice(Y.length), $.push(Y);
      } else throw Error("Unexpected dot segment condition");
    }
    return $.join("");
  }
  function L2(Q, $) {
    let X = $ !== true ? escape : unescape;
    if (Q.scheme !== void 0) Q.scheme = X(Q.scheme);
    if (Q.userinfo !== void 0) Q.userinfo = X(Q.userinfo);
    if (Q.host !== void 0) Q.host = X(Q.host);
    if (Q.path !== void 0) Q.path = X(Q.path);
    if (Q.query !== void 0) Q.query = X(Q.query);
    if (Q.fragment !== void 0) Q.fragment = X(Q.fragment);
    return Q;
  }
  function F2(Q) {
    let $ = [];
    if (Q.userinfo !== void 0) $.push(Q.userinfo), $.push("@");
    if (Q.host !== void 0) {
      let X = unescape(Q.host), Y = KB(X);
      if (Y.isIPV4) X = Y.host;
      else {
        let J = VB(Y.host);
        if (J.isIPV6 === true) X = `[${J.escapedHost}]`;
        else X = Q.host;
      }
      $.push(X);
    }
    if (typeof Q.port === "number" || typeof Q.port === "string") $.push(":"), $.push(String(Q.port));
    return $.length ? $.join("") : void 0;
  }
  UB.exports = { recomposeAuthority: F2, normalizeComponentEncoding: L2, removeDotSegments: U2, normalizeIPv4: KB, normalizeIPv6: VB, stringArrayToHexStripped: w7 };
});
var MB = P((Iy, wB) => {
  var N2 = /^[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}$/iu, O2 = /([\da-z][\d\-a-z]{0,31}):((?:[\w!$'()*+,\-.:;=@]|%[\da-f]{2})+)/iu;
  function FB(Q) {
    return typeof Q.secure === "boolean" ? Q.secure : String(Q.scheme).toLowerCase() === "wss";
  }
  function NB(Q) {
    if (!Q.host) Q.error = Q.error || "HTTP URIs must have a host.";
    return Q;
  }
  function OB(Q) {
    let $ = String(Q.scheme).toLowerCase() === "https";
    if (Q.port === ($ ? 443 : 80) || Q.port === "") Q.port = void 0;
    if (!Q.path) Q.path = "/";
    return Q;
  }
  function D2(Q) {
    return Q.secure = FB(Q), Q.resourceName = (Q.path || "/") + (Q.query ? "?" + Q.query : ""), Q.path = void 0, Q.query = void 0, Q;
  }
  function w2(Q) {
    if (Q.port === (FB(Q) ? 443 : 80) || Q.port === "") Q.port = void 0;
    if (typeof Q.secure === "boolean") Q.scheme = Q.secure ? "wss" : "ws", Q.secure = void 0;
    if (Q.resourceName) {
      let [$, X] = Q.resourceName.split("?");
      Q.path = $ && $ !== "/" ? $ : void 0, Q.query = X, Q.resourceName = void 0;
    }
    return Q.fragment = void 0, Q;
  }
  function M2(Q, $) {
    if (!Q.path) return Q.error = "URN can not be parsed", Q;
    let X = Q.path.match(O2);
    if (X) {
      let Y = $.scheme || Q.scheme || "urn";
      Q.nid = X[1].toLowerCase(), Q.nss = X[2];
      let J = `${Y}:${$.nid || Q.nid}`, G = M7[J];
      if (Q.path = void 0, G) Q = G.parse(Q, $);
    } else Q.error = Q.error || "URN can not be parsed.";
    return Q;
  }
  function A2(Q, $) {
    let X = $.scheme || Q.scheme || "urn", Y = Q.nid.toLowerCase(), J = `${X}:${$.nid || Y}`, G = M7[J];
    if (G) Q = G.serialize(Q, $);
    let W = Q, H = Q.nss;
    return W.path = `${Y || $.nid}:${H}`, $.skipEscape = true, W;
  }
  function j2(Q, $) {
    let X = Q;
    if (X.uuid = X.nss, X.nss = void 0, !$.tolerant && (!X.uuid || !N2.test(X.uuid))) X.error = X.error || "UUID is not valid.";
    return X;
  }
  function R2(Q) {
    let $ = Q;
    return $.nss = (Q.uuid || "").toLowerCase(), $;
  }
  var DB = { scheme: "http", domainHost: true, parse: NB, serialize: OB }, I2 = { scheme: "https", domainHost: DB.domainHost, parse: NB, serialize: OB }, OQ = { scheme: "ws", domainHost: true, parse: D2, serialize: w2 }, E2 = { scheme: "wss", domainHost: OQ.domainHost, parse: OQ.parse, serialize: OQ.serialize }, P2 = { scheme: "urn", parse: M2, serialize: A2, skipNormalize: true }, b2 = { scheme: "urn:uuid", parse: j2, serialize: R2, skipNormalize: true }, M7 = { http: DB, https: I2, ws: OQ, wss: E2, urn: P2, "urn:uuid": b2 };
  wB.exports = M7;
});
var jB = P((Ey, wQ) => {
  var { normalizeIPv6: Z2, normalizeIPv4: C2, removeDotSegments: S9, recomposeAuthority: S2, normalizeComponentEncoding: DQ } = LB(), A7 = MB();
  function _2(Q, $) {
    if (typeof Q === "string") Q = I1(T1(Q, $), $);
    else if (typeof Q === "object") Q = T1(I1(Q, $), $);
    return Q;
  }
  function k2(Q, $, X) {
    let Y = Object.assign({ scheme: "null" }, X), J = AB(T1(Q, Y), T1($, Y), Y, true);
    return I1(J, { ...Y, skipEscape: true });
  }
  function AB(Q, $, X, Y) {
    let J = {};
    if (!Y) Q = T1(I1(Q, X), X), $ = T1(I1($, X), X);
    if (X = X || {}, !X.tolerant && $.scheme) J.scheme = $.scheme, J.userinfo = $.userinfo, J.host = $.host, J.port = $.port, J.path = S9($.path || ""), J.query = $.query;
    else {
      if ($.userinfo !== void 0 || $.host !== void 0 || $.port !== void 0) J.userinfo = $.userinfo, J.host = $.host, J.port = $.port, J.path = S9($.path || ""), J.query = $.query;
      else {
        if (!$.path) if (J.path = Q.path, $.query !== void 0) J.query = $.query;
        else J.query = Q.query;
        else {
          if ($.path.charAt(0) === "/") J.path = S9($.path);
          else {
            if ((Q.userinfo !== void 0 || Q.host !== void 0 || Q.port !== void 0) && !Q.path) J.path = "/" + $.path;
            else if (!Q.path) J.path = $.path;
            else J.path = Q.path.slice(0, Q.path.lastIndexOf("/") + 1) + $.path;
            J.path = S9(J.path);
          }
          J.query = $.query;
        }
        J.userinfo = Q.userinfo, J.host = Q.host, J.port = Q.port;
      }
      J.scheme = Q.scheme;
    }
    return J.fragment = $.fragment, J;
  }
  function v2(Q, $, X) {
    if (typeof Q === "string") Q = unescape(Q), Q = I1(DQ(T1(Q, X), true), { ...X, skipEscape: true });
    else if (typeof Q === "object") Q = I1(DQ(Q, true), { ...X, skipEscape: true });
    if (typeof $ === "string") $ = unescape($), $ = I1(DQ(T1($, X), true), { ...X, skipEscape: true });
    else if (typeof $ === "object") $ = I1(DQ($, true), { ...X, skipEscape: true });
    return Q.toLowerCase() === $.toLowerCase();
  }
  function I1(Q, $) {
    let X = { host: Q.host, scheme: Q.scheme, userinfo: Q.userinfo, port: Q.port, path: Q.path, query: Q.query, nid: Q.nid, nss: Q.nss, uuid: Q.uuid, fragment: Q.fragment, reference: Q.reference, resourceName: Q.resourceName, secure: Q.secure, error: "" }, Y = Object.assign({}, $), J = [], G = A7[(Y.scheme || X.scheme || "").toLowerCase()];
    if (G && G.serialize) G.serialize(X, Y);
    if (X.path !== void 0) if (!Y.skipEscape) {
      if (X.path = escape(X.path), X.scheme !== void 0) X.path = X.path.split("%3A").join(":");
    } else X.path = unescape(X.path);
    if (Y.reference !== "suffix" && X.scheme) J.push(X.scheme, ":");
    let W = S2(X);
    if (W !== void 0) {
      if (Y.reference !== "suffix") J.push("//");
      if (J.push(W), X.path && X.path.charAt(0) !== "/") J.push("/");
    }
    if (X.path !== void 0) {
      let H = X.path;
      if (!Y.absolutePath && (!G || !G.absolutePath)) H = S9(H);
      if (W === void 0) H = H.replace(/^\/\//u, "/%2F");
      J.push(H);
    }
    if (X.query !== void 0) J.push("?", X.query);
    if (X.fragment !== void 0) J.push("#", X.fragment);
    return J.join("");
  }
  var T2 = Array.from({ length: 127 }, (Q, $) => /[^!"$&'()*+,\-.;=_`a-z{}~]/u.test(String.fromCharCode($)));
  function x2(Q) {
    let $ = 0;
    for (let X = 0, Y = Q.length; X < Y; ++X) if ($ = Q.charCodeAt(X), $ > 126 || T2[$]) return true;
    return false;
  }
  var y2 = /^(?:([^#/:?]+):)?(?:\/\/((?:([^#/?@]*)@)?(\[[^#/?\]]+\]|[^#/:?]*)(?::(\d*))?))?([^#?]*)(?:\?([^#]*))?(?:#((?:.|[\n\r])*))?/u;
  function T1(Q, $) {
    let X = Object.assign({}, $), Y = { scheme: void 0, userinfo: void 0, host: "", port: void 0, path: "", query: void 0, fragment: void 0 }, J = Q.indexOf("%") !== -1, G = false;
    if (X.reference === "suffix") Q = (X.scheme ? X.scheme + ":" : "") + "//" + Q;
    let W = Q.match(y2);
    if (W) {
      if (Y.scheme = W[1], Y.userinfo = W[3], Y.host = W[4], Y.port = parseInt(W[5], 10), Y.path = W[6] || "", Y.query = W[7], Y.fragment = W[8], isNaN(Y.port)) Y.port = W[5];
      if (Y.host) {
        let B = C2(Y.host);
        if (B.isIPV4 === false) {
          let z = Z2(B.host);
          Y.host = z.host.toLowerCase(), G = z.isIPV6;
        } else Y.host = B.host, G = true;
      }
      if (Y.scheme === void 0 && Y.userinfo === void 0 && Y.host === void 0 && Y.port === void 0 && Y.query === void 0 && !Y.path) Y.reference = "same-document";
      else if (Y.scheme === void 0) Y.reference = "relative";
      else if (Y.fragment === void 0) Y.reference = "absolute";
      else Y.reference = "uri";
      if (X.reference && X.reference !== "suffix" && X.reference !== Y.reference) Y.error = Y.error || "URI is not a " + X.reference + " reference.";
      let H = A7[(X.scheme || Y.scheme || "").toLowerCase()];
      if (!X.unicodeSupport && (!H || !H.unicodeSupport)) {
        if (Y.host && (X.domainHost || H && H.domainHost) && G === false && x2(Y.host)) try {
          Y.host = URL.domainToASCII(Y.host.toLowerCase());
        } catch (B) {
          Y.error = Y.error || "Host's domain name can not be converted to ASCII: " + B;
        }
      }
      if (!H || H && !H.skipNormalize) {
        if (J && Y.scheme !== void 0) Y.scheme = unescape(Y.scheme);
        if (J && Y.host !== void 0) Y.host = unescape(Y.host);
        if (Y.path) Y.path = escape(unescape(Y.path));
        if (Y.fragment) Y.fragment = encodeURI(decodeURIComponent(Y.fragment));
      }
      if (H && H.parse) H.parse(Y, X);
    } else Y.error = Y.error || "URI can not be parsed.";
    return Y;
  }
  var j7 = { SCHEMES: A7, normalize: _2, resolve: k2, resolveComponents: AB, equal: v2, serialize: I1, parse: T1 };
  wQ.exports = j7;
  wQ.exports.default = j7;
  wQ.exports.fastUri = j7;
});
var EB = P((IB) => {
  Object.defineProperty(IB, "__esModule", { value: true });
  var RB = jB();
  RB.code = 'require("ajv/dist/runtime/uri").default';
  IB.default = RB;
});
var vB = P((x1) => {
  Object.defineProperty(x1, "__esModule", { value: true });
  x1.CodeGen = x1.Name = x1.nil = x1.stringify = x1.str = x1._ = x1.KeywordCxt = void 0;
  var h2 = b9();
  Object.defineProperty(x1, "KeywordCxt", { enumerable: true, get: function() {
    return h2.KeywordCxt;
  } });
  var E4 = c();
  Object.defineProperty(x1, "_", { enumerable: true, get: function() {
    return E4._;
  } });
  Object.defineProperty(x1, "str", { enumerable: true, get: function() {
    return E4.str;
  } });
  Object.defineProperty(x1, "stringify", { enumerable: true, get: function() {
    return E4.stringify;
  } });
  Object.defineProperty(x1, "nil", { enumerable: true, get: function() {
    return E4.nil;
  } });
  Object.defineProperty(x1, "Name", { enumerable: true, get: function() {
    return E4.Name;
  } });
  Object.defineProperty(x1, "CodeGen", { enumerable: true, get: function() {
    return E4.CodeGen;
  } });
  var f2 = LQ(), SB = Z9(), u2 = X7(), _9 = NQ(), m2 = c(), k9 = I9(), MQ = R9(), I7 = a(), PB = JB(), l2 = EB(), _B = (Q, $) => new RegExp(Q, $);
  _B.code = "new RegExp";
  var c2 = ["removeAdditional", "useDefaults", "coerceTypes"], p2 = /* @__PURE__ */ new Set(["validate", "serialize", "parse", "wrapper", "root", "schema", "keyword", "pattern", "formats", "validate$data", "func", "obj", "Error"]), d2 = { errorDataPath: "", format: "`validateFormats: false` can be used instead.", nullable: '"nullable" keyword is supported by default.', jsonPointers: "Deprecated jsPropertySyntax can be used instead.", extendRefs: "Deprecated ignoreKeywordsWithRef can be used instead.", missingRefs: "Pass empty schema with $id that should be ignored to ajv.addSchema.", processCode: "Use option `code: {process: (code, schemaEnv: object) => string}`", sourceCode: "Use option `code: {source: true}`", strictDefaults: "It is default now, see option `strict`.", strictKeywords: "It is default now, see option `strict`.", uniqueItems: '"uniqueItems" keyword is always validated.', unknownFormats: "Disable strict mode or pass `true` to `ajv.addFormat` (or `formats` option).", cache: "Map is used as cache, schema object as key.", serialize: "Map is used as cache, schema object as key.", ajvErrors: "It is default now." }, i2 = { ignoreKeywordsWithRef: "", jsPropertySyntax: "", unicode: '"minLength"/"maxLength" account for unicode characters by default.' }, bB = 200;
  function n2(Q) {
    var $, X, Y, J, G, W, H, B, z, K, U, q, V, L, F, w, D, M, R, Z, v, O0, D0, d0, B6;
    let F1 = Q.strict, z6 = ($ = Q.code) === null || $ === void 0 ? void 0 : $.optimize, y1 = z6 === true || z6 === void 0 ? 1 : z6 || 0, K6 = (Y = (X = Q.code) === null || X === void 0 ? void 0 : X.regExp) !== null && Y !== void 0 ? Y : _B, h = (J = Q.uriResolver) !== null && J !== void 0 ? J : l2.default;
    return { strictSchema: (W = (G = Q.strictSchema) !== null && G !== void 0 ? G : F1) !== null && W !== void 0 ? W : true, strictNumbers: (B = (H = Q.strictNumbers) !== null && H !== void 0 ? H : F1) !== null && B !== void 0 ? B : true, strictTypes: (K = (z = Q.strictTypes) !== null && z !== void 0 ? z : F1) !== null && K !== void 0 ? K : "log", strictTuples: (q = (U = Q.strictTuples) !== null && U !== void 0 ? U : F1) !== null && q !== void 0 ? q : "log", strictRequired: (L = (V = Q.strictRequired) !== null && V !== void 0 ? V : F1) !== null && L !== void 0 ? L : false, code: Q.code ? { ...Q.code, optimize: y1, regExp: K6 } : { optimize: y1, regExp: K6 }, loopRequired: (F = Q.loopRequired) !== null && F !== void 0 ? F : bB, loopEnum: (w = Q.loopEnum) !== null && w !== void 0 ? w : bB, meta: (D = Q.meta) !== null && D !== void 0 ? D : true, messages: (M = Q.messages) !== null && M !== void 0 ? M : true, inlineRefs: (R = Q.inlineRefs) !== null && R !== void 0 ? R : true, schemaId: (Z = Q.schemaId) !== null && Z !== void 0 ? Z : "$id", addUsedSchema: (v = Q.addUsedSchema) !== null && v !== void 0 ? v : true, validateSchema: (O0 = Q.validateSchema) !== null && O0 !== void 0 ? O0 : true, validateFormats: (D0 = Q.validateFormats) !== null && D0 !== void 0 ? D0 : true, unicodeRegExp: (d0 = Q.unicodeRegExp) !== null && d0 !== void 0 ? d0 : true, int32range: (B6 = Q.int32range) !== null && B6 !== void 0 ? B6 : true, uriResolver: h };
  }
  class AQ {
    constructor(Q = {}) {
      this.schemas = {}, this.refs = {}, this.formats = {}, this._compilations = /* @__PURE__ */ new Set(), this._loading = {}, this._cache = /* @__PURE__ */ new Map(), Q = this.opts = { ...Q, ...n2(Q) };
      let { es5: $, lines: X } = this.opts.code;
      this.scope = new m2.ValueScope({ scope: {}, prefixes: p2, es5: $, lines: X }), this.logger = e2(Q.logger);
      let Y = Q.validateFormats;
      if (Q.validateFormats = false, this.RULES = (0, u2.getRules)(), ZB.call(this, d2, Q, "NOT SUPPORTED"), ZB.call(this, i2, Q, "DEPRECATED", "warn"), this._metaOpts = a2.call(this), Q.formats) r2.call(this);
      if (this._addVocabularies(), this._addDefaultMetaSchema(), Q.keywords) t2.call(this, Q.keywords);
      if (typeof Q.meta == "object") this.addMetaSchema(Q.meta);
      o2.call(this), Q.validateFormats = Y;
    }
    _addVocabularies() {
      this.addKeyword("$async");
    }
    _addDefaultMetaSchema() {
      let { $data: Q, meta: $, schemaId: X } = this.opts, Y = PB;
      if (X === "id") Y = { ...PB }, Y.id = Y.$id, delete Y.$id;
      if ($ && Q) this.addMetaSchema(Y, Y[X], false);
    }
    defaultMeta() {
      let { meta: Q, schemaId: $ } = this.opts;
      return this.opts.defaultMeta = typeof Q == "object" ? Q[$] || Q : void 0;
    }
    validate(Q, $) {
      let X;
      if (typeof Q == "string") {
        if (X = this.getSchema(Q), !X) throw Error(`no schema with key or ref "${Q}"`);
      } else X = this.compile(Q);
      let Y = X($);
      if (!("$async" in X)) this.errors = X.errors;
      return Y;
    }
    compile(Q, $) {
      let X = this._addSchema(Q, $);
      return X.validate || this._compileSchemaEnv(X);
    }
    compileAsync(Q, $) {
      if (typeof this.opts.loadSchema != "function") throw Error("options.loadSchema should be a function");
      let { loadSchema: X } = this.opts;
      return Y.call(this, Q, $);
      async function Y(z, K) {
        await J.call(this, z.$schema);
        let U = this._addSchema(z, K);
        return U.validate || G.call(this, U);
      }
      async function J(z) {
        if (z && !this.getSchema(z)) await Y.call(this, { $ref: z }, true);
      }
      async function G(z) {
        try {
          return this._compileSchemaEnv(z);
        } catch (K) {
          if (!(K instanceof SB.default)) throw K;
          return W.call(this, K), await H.call(this, K.missingSchema), G.call(this, z);
        }
      }
      function W({ missingSchema: z, missingRef: K }) {
        if (this.refs[z]) throw Error(`AnySchema ${z} is loaded but ${K} cannot be resolved`);
      }
      async function H(z) {
        let K = await B.call(this, z);
        if (!this.refs[z]) await J.call(this, K.$schema);
        if (!this.refs[z]) this.addSchema(K, z, $);
      }
      async function B(z) {
        let K = this._loading[z];
        if (K) return K;
        try {
          return await (this._loading[z] = X(z));
        } finally {
          delete this._loading[z];
        }
      }
    }
    addSchema(Q, $, X, Y = this.opts.validateSchema) {
      if (Array.isArray(Q)) {
        for (let G of Q) this.addSchema(G, void 0, X, Y);
        return this;
      }
      let J;
      if (typeof Q === "object") {
        let { schemaId: G } = this.opts;
        if (J = Q[G], J !== void 0 && typeof J != "string") throw Error(`schema ${G} must be string`);
      }
      return $ = (0, k9.normalizeId)($ || J), this._checkUnique($), this.schemas[$] = this._addSchema(Q, X, $, Y, true), this;
    }
    addMetaSchema(Q, $, X = this.opts.validateSchema) {
      return this.addSchema(Q, $, true, X), this;
    }
    validateSchema(Q, $) {
      if (typeof Q == "boolean") return true;
      let X;
      if (X = Q.$schema, X !== void 0 && typeof X != "string") throw Error("$schema must be a string");
      if (X = X || this.opts.defaultMeta || this.defaultMeta(), !X) return this.logger.warn("meta-schema not available"), this.errors = null, true;
      let Y = this.validate(X, Q);
      if (!Y && $) {
        let J = "schema is invalid: " + this.errorsText();
        if (this.opts.validateSchema === "log") this.logger.error(J);
        else throw Error(J);
      }
      return Y;
    }
    getSchema(Q) {
      let $;
      while (typeof ($ = CB.call(this, Q)) == "string") Q = $;
      if ($ === void 0) {
        let { schemaId: X } = this.opts, Y = new _9.SchemaEnv({ schema: {}, schemaId: X });
        if ($ = _9.resolveSchema.call(this, Y, Q), !$) return;
        this.refs[Q] = $;
      }
      return $.validate || this._compileSchemaEnv($);
    }
    removeSchema(Q) {
      if (Q instanceof RegExp) return this._removeAllSchemas(this.schemas, Q), this._removeAllSchemas(this.refs, Q), this;
      switch (typeof Q) {
        case "undefined":
          return this._removeAllSchemas(this.schemas), this._removeAllSchemas(this.refs), this._cache.clear(), this;
        case "string": {
          let $ = CB.call(this, Q);
          if (typeof $ == "object") this._cache.delete($.schema);
          return delete this.schemas[Q], delete this.refs[Q], this;
        }
        case "object": {
          let $ = Q;
          this._cache.delete($);
          let X = Q[this.opts.schemaId];
          if (X) X = (0, k9.normalizeId)(X), delete this.schemas[X], delete this.refs[X];
          return this;
        }
        default:
          throw Error("ajv.removeSchema: invalid parameter");
      }
    }
    addVocabulary(Q) {
      for (let $ of Q) this.addKeyword($);
      return this;
    }
    addKeyword(Q, $) {
      let X;
      if (typeof Q == "string") {
        if (X = Q, typeof $ == "object") this.logger.warn("these parameters are deprecated, see docs for addKeyword"), $.keyword = X;
      } else if (typeof Q == "object" && $ === void 0) {
        if ($ = Q, X = $.keyword, Array.isArray(X) && !X.length) throw Error("addKeywords: keyword must be string or non-empty array");
      } else throw Error("invalid addKeywords parameters");
      if ($R.call(this, X, $), !$) return (0, I7.eachItem)(X, (J) => R7.call(this, J)), this;
      YR.call(this, $);
      let Y = { ...$, type: (0, MQ.getJSONTypes)($.type), schemaType: (0, MQ.getJSONTypes)($.schemaType) };
      return (0, I7.eachItem)(X, Y.type.length === 0 ? (J) => R7.call(this, J, Y) : (J) => Y.type.forEach((G) => R7.call(this, J, Y, G))), this;
    }
    getKeyword(Q) {
      let $ = this.RULES.all[Q];
      return typeof $ == "object" ? $.definition : !!$;
    }
    removeKeyword(Q) {
      let { RULES: $ } = this;
      delete $.keywords[Q], delete $.all[Q];
      for (let X of $.rules) {
        let Y = X.rules.findIndex((J) => J.keyword === Q);
        if (Y >= 0) X.rules.splice(Y, 1);
      }
      return this;
    }
    addFormat(Q, $) {
      if (typeof $ == "string") $ = new RegExp($);
      return this.formats[Q] = $, this;
    }
    errorsText(Q = this.errors, { separator: $ = ", ", dataVar: X = "data" } = {}) {
      if (!Q || Q.length === 0) return "No errors";
      return Q.map((Y) => `${X}${Y.instancePath} ${Y.message}`).reduce((Y, J) => Y + $ + J);
    }
    $dataMetaSchema(Q, $) {
      let X = this.RULES.all;
      Q = JSON.parse(JSON.stringify(Q));
      for (let Y of $) {
        let J = Y.split("/").slice(1), G = Q;
        for (let W of J) G = G[W];
        for (let W in X) {
          let H = X[W];
          if (typeof H != "object") continue;
          let { $data: B } = H.definition, z = G[W];
          if (B && z) G[W] = kB(z);
        }
      }
      return Q;
    }
    _removeAllSchemas(Q, $) {
      for (let X in Q) {
        let Y = Q[X];
        if (!$ || $.test(X)) {
          if (typeof Y == "string") delete Q[X];
          else if (Y && !Y.meta) this._cache.delete(Y.schema), delete Q[X];
        }
      }
    }
    _addSchema(Q, $, X, Y = this.opts.validateSchema, J = this.opts.addUsedSchema) {
      let G, { schemaId: W } = this.opts;
      if (typeof Q == "object") G = Q[W];
      else if (this.opts.jtd) throw Error("schema must be object");
      else if (typeof Q != "boolean") throw Error("schema must be object or boolean");
      let H = this._cache.get(Q);
      if (H !== void 0) return H;
      X = (0, k9.normalizeId)(G || X);
      let B = k9.getSchemaRefs.call(this, Q, X);
      if (H = new _9.SchemaEnv({ schema: Q, schemaId: W, meta: $, baseId: X, localRefs: B }), this._cache.set(H.schema, H), J && !X.startsWith("#")) {
        if (X) this._checkUnique(X);
        this.refs[X] = H;
      }
      if (Y) this.validateSchema(Q, true);
      return H;
    }
    _checkUnique(Q) {
      if (this.schemas[Q] || this.refs[Q]) throw Error(`schema with key or id "${Q}" already exists`);
    }
    _compileSchemaEnv(Q) {
      if (Q.meta) this._compileMetaSchema(Q);
      else _9.compileSchema.call(this, Q);
      if (!Q.validate) throw Error("ajv implementation error");
      return Q.validate;
    }
    _compileMetaSchema(Q) {
      let $ = this.opts;
      this.opts = this._metaOpts;
      try {
        _9.compileSchema.call(this, Q);
      } finally {
        this.opts = $;
      }
    }
  }
  AQ.ValidationError = f2.default;
  AQ.MissingRefError = SB.default;
  x1.default = AQ;
  function ZB(Q, $, X, Y = "error") {
    for (let J in Q) {
      let G = J;
      if (G in $) this.logger[Y](`${X}: option ${J}. ${Q[G]}`);
    }
  }
  function CB(Q) {
    return Q = (0, k9.normalizeId)(Q), this.schemas[Q] || this.refs[Q];
  }
  function o2() {
    let Q = this.opts.schemas;
    if (!Q) return;
    if (Array.isArray(Q)) this.addSchema(Q);
    else for (let $ in Q) this.addSchema(Q[$], $);
  }
  function r2() {
    for (let Q in this.opts.formats) {
      let $ = this.opts.formats[Q];
      if ($) this.addFormat(Q, $);
    }
  }
  function t2(Q) {
    if (Array.isArray(Q)) {
      this.addVocabulary(Q);
      return;
    }
    this.logger.warn("keywords option as map is deprecated, pass array");
    for (let $ in Q) {
      let X = Q[$];
      if (!X.keyword) X.keyword = $;
      this.addKeyword(X);
    }
  }
  function a2() {
    let Q = { ...this.opts };
    for (let $ of c2) delete Q[$];
    return Q;
  }
  var s2 = { log() {
  }, warn() {
  }, error() {
  } };
  function e2(Q) {
    if (Q === false) return s2;
    if (Q === void 0) return console;
    if (Q.log && Q.warn && Q.error) return Q;
    throw Error("logger must implement log, warn and error methods");
  }
  var QR = /^[a-z_$][a-z0-9_$:-]*$/i;
  function $R(Q, $) {
    let { RULES: X } = this;
    if ((0, I7.eachItem)(Q, (Y) => {
      if (X.keywords[Y]) throw Error(`Keyword ${Y} is already defined`);
      if (!QR.test(Y)) throw Error(`Keyword ${Y} has invalid name`);
    }), !$) return;
    if ($.$data && !("code" in $ || "validate" in $)) throw Error('$data keyword must have "code" or "validate" function');
  }
  function R7(Q, $, X) {
    var Y;
    let J = $ === null || $ === void 0 ? void 0 : $.post;
    if (X && J) throw Error('keyword with "post" flag cannot have "type"');
    let { RULES: G } = this, W = J ? G.post : G.rules.find(({ type: B }) => B === X);
    if (!W) W = { type: X, rules: [] }, G.rules.push(W);
    if (G.keywords[Q] = true, !$) return;
    let H = { keyword: Q, definition: { ...$, type: (0, MQ.getJSONTypes)($.type), schemaType: (0, MQ.getJSONTypes)($.schemaType) } };
    if ($.before) XR.call(this, W, H, $.before);
    else W.rules.push(H);
    G.all[Q] = H, (Y = $.implements) === null || Y === void 0 || Y.forEach((B) => this.addKeyword(B));
  }
  function XR(Q, $, X) {
    let Y = Q.rules.findIndex((J) => J.keyword === X);
    if (Y >= 0) Q.rules.splice(Y, 0, $);
    else Q.rules.push($), this.logger.warn(`rule ${X} is not defined`);
  }
  function YR(Q) {
    let { metaSchema: $ } = Q;
    if ($ === void 0) return;
    if (Q.$data && this.opts.$data) $ = kB($);
    Q.validateSchema = this.compile($, true);
  }
  var JR = { $ref: "https://raw.githubusercontent.com/ajv-validator/ajv/master/lib/refs/data.json#" };
  function kB(Q) {
    return { anyOf: [Q, JR] };
  }
});
var xB = P((TB) => {
  Object.defineProperty(TB, "__esModule", { value: true });
  var HR = { keyword: "id", code() {
    throw Error('NOT SUPPORTED: keyword "id", use "$id" for schema ID');
  } };
  TB.default = HR;
});
var mB = P((fB) => {
  Object.defineProperty(fB, "__esModule", { value: true });
  fB.callRef = fB.getValidate = void 0;
  var zR = Z9(), yB = e0(), f0 = c(), P4 = k1(), gB = NQ(), jQ = a(), KR = { keyword: "$ref", schemaType: "string", code(Q) {
    let { gen: $, schema: X, it: Y } = Q, { baseId: J, schemaEnv: G, validateName: W, opts: H, self: B } = Y, { root: z } = G;
    if ((X === "#" || X === "#/") && J === z.baseId) return U();
    let K = gB.resolveRef.call(B, z, J, X);
    if (K === void 0) throw new zR.default(Y.opts.uriResolver, J, X);
    if (K instanceof gB.SchemaEnv) return q(K);
    return V(K);
    function U() {
      if (G === z) return RQ(Q, W, G, G.$async);
      let L = $.scopeValue("root", { ref: z });
      return RQ(Q, f0._`${L}.validate`, z, z.$async);
    }
    function q(L) {
      let F = hB(Q, L);
      RQ(Q, F, L, L.$async);
    }
    function V(L) {
      let F = $.scopeValue("schema", H.code.source === true ? { ref: L, code: (0, f0.stringify)(L) } : { ref: L }), w = $.name("valid"), D = Q.subschema({ schema: L, dataTypes: [], schemaPath: f0.nil, topSchemaRef: F, errSchemaPath: X }, w);
      Q.mergeEvaluated(D), Q.ok(w);
    }
  } };
  function hB(Q, $) {
    let { gen: X } = Q;
    return $.validate ? X.scopeValue("validate", { ref: $.validate }) : f0._`${X.scopeValue("wrapper", { ref: $ })}.validate`;
  }
  fB.getValidate = hB;
  function RQ(Q, $, X, Y) {
    let { gen: J, it: G } = Q, { allErrors: W, schemaEnv: H, opts: B } = G, z = B.passContext ? P4.default.this : f0.nil;
    if (Y) K();
    else U();
    function K() {
      if (!H.$async) throw Error("async schema referenced by sync schema");
      let L = J.let("valid");
      J.try(() => {
        if (J.code(f0._`await ${(0, yB.callValidateCode)(Q, $, z)}`), V($), !W) J.assign(L, true);
      }, (F) => {
        if (J.if(f0._`!(${F} instanceof ${G.ValidationError})`, () => J.throw(F)), q(F), !W) J.assign(L, false);
      }), Q.ok(L);
    }
    function U() {
      Q.result((0, yB.callValidateCode)(Q, $, z), () => V($), () => q($));
    }
    function q(L) {
      let F = f0._`${L}.errors`;
      J.assign(P4.default.vErrors, f0._`${P4.default.vErrors} === null ? ${F} : ${P4.default.vErrors}.concat(${F})`), J.assign(P4.default.errors, f0._`${P4.default.vErrors}.length`);
    }
    function V(L) {
      var F;
      if (!G.opts.unevaluated) return;
      let w = (F = X === null || X === void 0 ? void 0 : X.validate) === null || F === void 0 ? void 0 : F.evaluated;
      if (G.props !== true) if (w && !w.dynamicProps) {
        if (w.props !== void 0) G.props = jQ.mergeEvaluated.props(J, w.props, G.props);
      } else {
        let D = J.var("props", f0._`${L}.evaluated.props`);
        G.props = jQ.mergeEvaluated.props(J, D, G.props, f0.Name);
      }
      if (G.items !== true) if (w && !w.dynamicItems) {
        if (w.items !== void 0) G.items = jQ.mergeEvaluated.items(J, w.items, G.items);
      } else {
        let D = J.var("items", f0._`${L}.evaluated.items`);
        G.items = jQ.mergeEvaluated.items(J, D, G.items, f0.Name);
      }
    }
  }
  fB.callRef = RQ;
  fB.default = KR;
});
var cB = P((lB) => {
  Object.defineProperty(lB, "__esModule", { value: true });
  var UR = xB(), LR = mB(), FR = ["$schema", "$id", "$defs", "$vocabulary", { keyword: "$comment" }, "definitions", UR.default, LR.default];
  lB.default = FR;
});
var dB = P((pB) => {
  Object.defineProperty(pB, "__esModule", { value: true });
  var IQ = c(), G6 = IQ.operators, EQ = { maximum: { okStr: "<=", ok: G6.LTE, fail: G6.GT }, minimum: { okStr: ">=", ok: G6.GTE, fail: G6.LT }, exclusiveMaximum: { okStr: "<", ok: G6.LT, fail: G6.GTE }, exclusiveMinimum: { okStr: ">", ok: G6.GT, fail: G6.LTE } }, OR = { message: ({ keyword: Q, schemaCode: $ }) => IQ.str`must be ${EQ[Q].okStr} ${$}`, params: ({ keyword: Q, schemaCode: $ }) => IQ._`{comparison: ${EQ[Q].okStr}, limit: ${$}}` }, DR = { keyword: Object.keys(EQ), type: "number", schemaType: "number", $data: true, error: OR, code(Q) {
    let { keyword: $, data: X, schemaCode: Y } = Q;
    Q.fail$data(IQ._`${X} ${EQ[$].fail} ${Y} || isNaN(${X})`);
  } };
  pB.default = DR;
});
var nB = P((iB) => {
  Object.defineProperty(iB, "__esModule", { value: true });
  var v9 = c(), MR = { message: ({ schemaCode: Q }) => v9.str`must be multiple of ${Q}`, params: ({ schemaCode: Q }) => v9._`{multipleOf: ${Q}}` }, AR = { keyword: "multipleOf", type: "number", schemaType: "number", $data: true, error: MR, code(Q) {
    let { gen: $, data: X, schemaCode: Y, it: J } = Q, G = J.opts.multipleOfPrecision, W = $.let("res"), H = G ? v9._`Math.abs(Math.round(${W}) - ${W}) > 1e-${G}` : v9._`${W} !== parseInt(${W})`;
    Q.fail$data(v9._`(${Y} === 0 || (${W} = ${X}/${Y}, ${H}))`);
  } };
  iB.default = AR;
});
var tB = P((rB) => {
  Object.defineProperty(rB, "__esModule", { value: true });
  function oB(Q) {
    let $ = Q.length, X = 0, Y = 0, J;
    while (Y < $) if (X++, J = Q.charCodeAt(Y++), J >= 55296 && J <= 56319 && Y < $) {
      if (J = Q.charCodeAt(Y), (J & 64512) === 56320) Y++;
    }
    return X;
  }
  rB.default = oB;
  oB.code = 'require("ajv/dist/runtime/ucs2length").default';
});
var sB = P((aB) => {
  Object.defineProperty(aB, "__esModule", { value: true });
  var _6 = c(), IR = a(), ER = tB(), PR = { message({ keyword: Q, schemaCode: $ }) {
    let X = Q === "maxLength" ? "more" : "fewer";
    return _6.str`must NOT have ${X} than ${$} characters`;
  }, params: ({ schemaCode: Q }) => _6._`{limit: ${Q}}` }, bR = { keyword: ["maxLength", "minLength"], type: "string", schemaType: "number", $data: true, error: PR, code(Q) {
    let { keyword: $, data: X, schemaCode: Y, it: J } = Q, G = $ === "maxLength" ? _6.operators.GT : _6.operators.LT, W = J.opts.unicode === false ? _6._`${X}.length` : _6._`${(0, IR.useFunc)(Q.gen, ER.default)}(${X})`;
    Q.fail$data(_6._`${W} ${G} ${Y}`);
  } };
  aB.default = bR;
});
var Qz = P((eB) => {
  Object.defineProperty(eB, "__esModule", { value: true });
  var CR = e0(), PQ = c(), SR = { message: ({ schemaCode: Q }) => PQ.str`must match pattern "${Q}"`, params: ({ schemaCode: Q }) => PQ._`{pattern: ${Q}}` }, _R = { keyword: "pattern", type: "string", schemaType: "string", $data: true, error: SR, code(Q) {
    let { data: $, $data: X, schema: Y, schemaCode: J, it: G } = Q, W = G.opts.unicodeRegExp ? "u" : "", H = X ? PQ._`(new RegExp(${J}, ${W}))` : (0, CR.usePattern)(Q, Y);
    Q.fail$data(PQ._`!${H}.test(${$})`);
  } };
  eB.default = _R;
});
var Xz = P(($z) => {
  Object.defineProperty($z, "__esModule", { value: true });
  var T9 = c(), vR = { message({ keyword: Q, schemaCode: $ }) {
    let X = Q === "maxProperties" ? "more" : "fewer";
    return T9.str`must NOT have ${X} than ${$} properties`;
  }, params: ({ schemaCode: Q }) => T9._`{limit: ${Q}}` }, TR = { keyword: ["maxProperties", "minProperties"], type: "object", schemaType: "number", $data: true, error: vR, code(Q) {
    let { keyword: $, data: X, schemaCode: Y } = Q, J = $ === "maxProperties" ? T9.operators.GT : T9.operators.LT;
    Q.fail$data(T9._`Object.keys(${X}).length ${J} ${Y}`);
  } };
  $z.default = TR;
});
var Jz = P((Yz) => {
  Object.defineProperty(Yz, "__esModule", { value: true });
  var x9 = e0(), y9 = c(), yR = a(), gR = { message: ({ params: { missingProperty: Q } }) => y9.str`must have required property '${Q}'`, params: ({ params: { missingProperty: Q } }) => y9._`{missingProperty: ${Q}}` }, hR = { keyword: "required", type: "object", schemaType: "array", $data: true, error: gR, code(Q) {
    let { gen: $, schema: X, schemaCode: Y, data: J, $data: G, it: W } = Q, { opts: H } = W;
    if (!G && X.length === 0) return;
    let B = X.length >= H.loopRequired;
    if (W.allErrors) z();
    else K();
    if (H.strictRequired) {
      let V = Q.parentSchema.properties, { definedProperties: L } = Q.it;
      for (let F of X) if ((V === null || V === void 0 ? void 0 : V[F]) === void 0 && !L.has(F)) {
        let w = W.schemaEnv.baseId + W.errSchemaPath, D = `required property "${F}" is not defined at "${w}" (strictRequired)`;
        (0, yR.checkStrictMode)(W, D, W.opts.strictRequired);
      }
    }
    function z() {
      if (B || G) Q.block$data(y9.nil, U);
      else for (let V of X) (0, x9.checkReportMissingProp)(Q, V);
    }
    function K() {
      let V = $.let("missing");
      if (B || G) {
        let L = $.let("valid", true);
        Q.block$data(L, () => q(V, L)), Q.ok(L);
      } else $.if((0, x9.checkMissingProp)(Q, X, V)), (0, x9.reportMissingProp)(Q, V), $.else();
    }
    function U() {
      $.forOf("prop", Y, (V) => {
        Q.setParams({ missingProperty: V }), $.if((0, x9.noPropertyInData)($, J, V, H.ownProperties), () => Q.error());
      });
    }
    function q(V, L) {
      Q.setParams({ missingProperty: V }), $.forOf(V, Y, () => {
        $.assign(L, (0, x9.propertyInData)($, J, V, H.ownProperties)), $.if((0, y9.not)(L), () => {
          Q.error(), $.break();
        });
      }, y9.nil);
    }
  } };
  Yz.default = hR;
});
var Wz = P((Gz) => {
  Object.defineProperty(Gz, "__esModule", { value: true });
  var g9 = c(), uR = { message({ keyword: Q, schemaCode: $ }) {
    let X = Q === "maxItems" ? "more" : "fewer";
    return g9.str`must NOT have ${X} than ${$} items`;
  }, params: ({ schemaCode: Q }) => g9._`{limit: ${Q}}` }, mR = { keyword: ["maxItems", "minItems"], type: "array", schemaType: "number", $data: true, error: uR, code(Q) {
    let { keyword: $, data: X, schemaCode: Y } = Q, J = $ === "maxItems" ? g9.operators.GT : g9.operators.LT;
    Q.fail$data(g9._`${X}.length ${J} ${Y}`);
  } };
  Gz.default = mR;
});
var bQ = P((Bz) => {
  Object.defineProperty(Bz, "__esModule", { value: true });
  var Hz = K7();
  Hz.code = 'require("ajv/dist/runtime/equal").default';
  Bz.default = Hz;
});
var Kz = P((zz) => {
  Object.defineProperty(zz, "__esModule", { value: true });
  var E7 = R9(), P0 = c(), pR = a(), dR = bQ(), iR = { message: ({ params: { i: Q, j: $ } }) => P0.str`must NOT have duplicate items (items ## ${$} and ${Q} are identical)`, params: ({ params: { i: Q, j: $ } }) => P0._`{i: ${Q}, j: ${$}}` }, nR = { keyword: "uniqueItems", type: "array", schemaType: "boolean", $data: true, error: iR, code(Q) {
    let { gen: $, data: X, $data: Y, schema: J, parentSchema: G, schemaCode: W, it: H } = Q;
    if (!Y && !J) return;
    let B = $.let("valid"), z = G.items ? (0, E7.getSchemaTypes)(G.items) : [];
    Q.block$data(B, K, P0._`${W} === false`), Q.ok(B);
    function K() {
      let L = $.let("i", P0._`${X}.length`), F = $.let("j");
      Q.setParams({ i: L, j: F }), $.assign(B, true), $.if(P0._`${L} > 1`, () => (U() ? q : V)(L, F));
    }
    function U() {
      return z.length > 0 && !z.some((L) => L === "object" || L === "array");
    }
    function q(L, F) {
      let w = $.name("item"), D = (0, E7.checkDataTypes)(z, w, H.opts.strictNumbers, E7.DataType.Wrong), M = $.const("indices", P0._`{}`);
      $.for(P0._`;${L}--;`, () => {
        if ($.let(w, P0._`${X}[${L}]`), $.if(D, P0._`continue`), z.length > 1) $.if(P0._`typeof ${w} == "string"`, P0._`${w} += "_"`);
        $.if(P0._`typeof ${M}[${w}] == "number"`, () => {
          $.assign(F, P0._`${M}[${w}]`), Q.error(), $.assign(B, false).break();
        }).code(P0._`${M}[${w}] = ${L}`);
      });
    }
    function V(L, F) {
      let w = (0, pR.useFunc)($, dR.default), D = $.name("outer");
      $.label(D).for(P0._`;${L}--;`, () => $.for(P0._`${F} = ${L}; ${F}--;`, () => $.if(P0._`${w}(${X}[${L}], ${X}[${F}])`, () => {
        Q.error(), $.assign(B, false).break(D);
      })));
    }
  } };
  zz.default = nR;
});
var qz = P((Vz) => {
  Object.defineProperty(Vz, "__esModule", { value: true });
  var P7 = c(), rR = a(), tR = bQ(), aR = { message: "must be equal to constant", params: ({ schemaCode: Q }) => P7._`{allowedValue: ${Q}}` }, sR = { keyword: "const", $data: true, error: aR, code(Q) {
    let { gen: $, data: X, $data: Y, schemaCode: J, schema: G } = Q;
    if (Y || G && typeof G == "object") Q.fail$data(P7._`!${(0, rR.useFunc)($, tR.default)}(${X}, ${J})`);
    else Q.fail(P7._`${G} !== ${X}`);
  } };
  Vz.default = sR;
});
var Lz = P((Uz) => {
  Object.defineProperty(Uz, "__esModule", { value: true });
  var h9 = c(), QI = a(), $I = bQ(), XI = { message: "must be equal to one of the allowed values", params: ({ schemaCode: Q }) => h9._`{allowedValues: ${Q}}` }, YI = { keyword: "enum", schemaType: "array", $data: true, error: XI, code(Q) {
    let { gen: $, data: X, $data: Y, schema: J, schemaCode: G, it: W } = Q;
    if (!Y && J.length === 0) throw Error("enum must have non-empty array");
    let H = J.length >= W.opts.loopEnum, B, z = () => B !== null && B !== void 0 ? B : B = (0, QI.useFunc)($, $I.default), K;
    if (H || Y) K = $.let("valid"), Q.block$data(K, U);
    else {
      if (!Array.isArray(J)) throw Error("ajv implementation error");
      let V = $.const("vSchema", G);
      K = (0, h9.or)(...J.map((L, F) => q(V, F)));
    }
    Q.pass(K);
    function U() {
      $.assign(K, false), $.forOf("v", G, (V) => $.if(h9._`${z()}(${X}, ${V})`, () => $.assign(K, true).break()));
    }
    function q(V, L) {
      let F = J[L];
      return typeof F === "object" && F !== null ? h9._`${z()}(${X}, ${V}[${L}])` : h9._`${X} === ${F}`;
    }
  } };
  Uz.default = YI;
});
var Nz = P((Fz) => {
  Object.defineProperty(Fz, "__esModule", { value: true });
  var GI = dB(), WI = nB(), HI = sB(), BI = Qz(), zI = Xz(), KI = Jz(), VI = Wz(), qI = Kz(), UI = qz(), LI = Lz(), FI = [GI.default, WI.default, HI.default, BI.default, zI.default, KI.default, VI.default, qI.default, { keyword: "type", schemaType: ["string", "array"] }, { keyword: "nullable", schemaType: "boolean" }, UI.default, LI.default];
  Fz.default = FI;
});
var Z7 = P((Dz) => {
  Object.defineProperty(Dz, "__esModule", { value: true });
  Dz.validateAdditionalItems = void 0;
  var k6 = c(), b7 = a(), OI = { message: ({ params: { len: Q } }) => k6.str`must NOT have more than ${Q} items`, params: ({ params: { len: Q } }) => k6._`{limit: ${Q}}` }, DI = { keyword: "additionalItems", type: "array", schemaType: ["boolean", "object"], before: "uniqueItems", error: OI, code(Q) {
    let { parentSchema: $, it: X } = Q, { items: Y } = $;
    if (!Array.isArray(Y)) {
      (0, b7.checkStrictMode)(X, '"additionalItems" is ignored when "items" is not an array of schemas');
      return;
    }
    Oz(Q, Y);
  } };
  function Oz(Q, $) {
    let { gen: X, schema: Y, data: J, keyword: G, it: W } = Q;
    W.items = true;
    let H = X.const("len", k6._`${J}.length`);
    if (Y === false) Q.setParams({ len: $.length }), Q.pass(k6._`${H} <= ${$.length}`);
    else if (typeof Y == "object" && !(0, b7.alwaysValidSchema)(W, Y)) {
      let z = X.var("valid", k6._`${H} <= ${$.length}`);
      X.if((0, k6.not)(z), () => B(z)), Q.ok(z);
    }
    function B(z) {
      X.forRange("i", $.length, H, (K) => {
        if (Q.subschema({ keyword: G, dataProp: K, dataPropType: b7.Type.Num }, z), !W.allErrors) X.if((0, k6.not)(z), () => X.break());
      });
    }
  }
  Dz.validateAdditionalItems = Oz;
  Dz.default = DI;
});
var C7 = P((jz) => {
  Object.defineProperty(jz, "__esModule", { value: true });
  jz.validateTuple = void 0;
  var Mz = c(), ZQ = a(), MI = e0(), AI = { keyword: "items", type: "array", schemaType: ["object", "array", "boolean"], before: "uniqueItems", code(Q) {
    let { schema: $, it: X } = Q;
    if (Array.isArray($)) return Az(Q, "additionalItems", $);
    if (X.items = true, (0, ZQ.alwaysValidSchema)(X, $)) return;
    Q.ok((0, MI.validateArray)(Q));
  } };
  function Az(Q, $, X = Q.schema) {
    let { gen: Y, parentSchema: J, data: G, keyword: W, it: H } = Q;
    if (K(J), H.opts.unevaluated && X.length && H.items !== true) H.items = ZQ.mergeEvaluated.items(Y, X.length, H.items);
    let B = Y.name("valid"), z = Y.const("len", Mz._`${G}.length`);
    X.forEach((U, q) => {
      if ((0, ZQ.alwaysValidSchema)(H, U)) return;
      Y.if(Mz._`${z} > ${q}`, () => Q.subschema({ keyword: W, schemaProp: q, dataProp: q }, B)), Q.ok(B);
    });
    function K(U) {
      let { opts: q, errSchemaPath: V } = H, L = X.length, F = L === U.minItems && (L === U.maxItems || U[$] === false);
      if (q.strictTuples && !F) {
        let w = `"${W}" is ${L}-tuple, but minItems or maxItems/${$} are not specified or different at path "${V}"`;
        (0, ZQ.checkStrictMode)(H, w, q.strictTuples);
      }
    }
  }
  jz.validateTuple = Az;
  jz.default = AI;
});
var Ez = P((Iz) => {
  Object.defineProperty(Iz, "__esModule", { value: true });
  var RI = C7(), II = { keyword: "prefixItems", type: "array", schemaType: ["array"], before: "uniqueItems", code: (Q) => (0, RI.validateTuple)(Q, "items") };
  Iz.default = II;
});
var Zz = P((bz) => {
  Object.defineProperty(bz, "__esModule", { value: true });
  var Pz = c(), PI = a(), bI = e0(), ZI = Z7(), CI = { message: ({ params: { len: Q } }) => Pz.str`must NOT have more than ${Q} items`, params: ({ params: { len: Q } }) => Pz._`{limit: ${Q}}` }, SI = { keyword: "items", type: "array", schemaType: ["object", "boolean"], before: "uniqueItems", error: CI, code(Q) {
    let { schema: $, parentSchema: X, it: Y } = Q, { prefixItems: J } = X;
    if (Y.items = true, (0, PI.alwaysValidSchema)(Y, $)) return;
    if (J) (0, ZI.validateAdditionalItems)(Q, J);
    else Q.ok((0, bI.validateArray)(Q));
  } };
  bz.default = SI;
});
var Sz = P((Cz) => {
  Object.defineProperty(Cz, "__esModule", { value: true });
  var Q1 = c(), CQ = a(), kI = { message: ({ params: { min: Q, max: $ } }) => $ === void 0 ? Q1.str`must contain at least ${Q} valid item(s)` : Q1.str`must contain at least ${Q} and no more than ${$} valid item(s)`, params: ({ params: { min: Q, max: $ } }) => $ === void 0 ? Q1._`{minContains: ${Q}}` : Q1._`{minContains: ${Q}, maxContains: ${$}}` }, vI = { keyword: "contains", type: "array", schemaType: ["object", "boolean"], before: "uniqueItems", trackErrors: true, error: kI, code(Q) {
    let { gen: $, schema: X, parentSchema: Y, data: J, it: G } = Q, W, H, { minContains: B, maxContains: z } = Y;
    if (G.opts.next) W = B === void 0 ? 1 : B, H = z;
    else W = 1;
    let K = $.const("len", Q1._`${J}.length`);
    if (Q.setParams({ min: W, max: H }), H === void 0 && W === 0) {
      (0, CQ.checkStrictMode)(G, '"minContains" == 0 without "maxContains": "contains" keyword ignored');
      return;
    }
    if (H !== void 0 && W > H) {
      (0, CQ.checkStrictMode)(G, '"minContains" > "maxContains" is always invalid'), Q.fail();
      return;
    }
    if ((0, CQ.alwaysValidSchema)(G, X)) {
      let F = Q1._`${K} >= ${W}`;
      if (H !== void 0) F = Q1._`${F} && ${K} <= ${H}`;
      Q.pass(F);
      return;
    }
    G.items = true;
    let U = $.name("valid");
    if (H === void 0 && W === 1) V(U, () => $.if(U, () => $.break()));
    else if (W === 0) {
      if ($.let(U, true), H !== void 0) $.if(Q1._`${J}.length > 0`, q);
    } else $.let(U, false), q();
    Q.result(U, () => Q.reset());
    function q() {
      let F = $.name("_valid"), w = $.let("count", 0);
      V(F, () => $.if(F, () => L(w)));
    }
    function V(F, w) {
      $.forRange("i", 0, K, (D) => {
        Q.subschema({ keyword: "contains", dataProp: D, dataPropType: CQ.Type.Num, compositeRule: true }, F), w();
      });
    }
    function L(F) {
      if ($.code(Q1._`${F}++`), H === void 0) $.if(Q1._`${F} >= ${W}`, () => $.assign(U, true).break());
      else if ($.if(Q1._`${F} > ${H}`, () => $.assign(U, false).break()), W === 1) $.assign(U, true);
      else $.if(Q1._`${F} >= ${W}`, () => $.assign(U, true));
    }
  } };
  Cz.default = vI;
});
var yz = P((vz) => {
  Object.defineProperty(vz, "__esModule", { value: true });
  vz.validateSchemaDeps = vz.validatePropertyDeps = vz.error = void 0;
  var S7 = c(), xI = a(), f9 = e0();
  vz.error = { message: ({ params: { property: Q, depsCount: $, deps: X } }) => {
    let Y = $ === 1 ? "property" : "properties";
    return S7.str`must have ${Y} ${X} when property ${Q} is present`;
  }, params: ({ params: { property: Q, depsCount: $, deps: X, missingProperty: Y } }) => S7._`{property: ${Q},
    missingProperty: ${Y},
    depsCount: ${$},
    deps: ${X}}` };
  var yI = { keyword: "dependencies", type: "object", schemaType: "object", error: vz.error, code(Q) {
    let [$, X] = gI(Q);
    _z(Q, $), kz(Q, X);
  } };
  function gI({ schema: Q }) {
    let $ = {}, X = {};
    for (let Y in Q) {
      if (Y === "__proto__") continue;
      let J = Array.isArray(Q[Y]) ? $ : X;
      J[Y] = Q[Y];
    }
    return [$, X];
  }
  function _z(Q, $ = Q.schema) {
    let { gen: X, data: Y, it: J } = Q;
    if (Object.keys($).length === 0) return;
    let G = X.let("missing");
    for (let W in $) {
      let H = $[W];
      if (H.length === 0) continue;
      let B = (0, f9.propertyInData)(X, Y, W, J.opts.ownProperties);
      if (Q.setParams({ property: W, depsCount: H.length, deps: H.join(", ") }), J.allErrors) X.if(B, () => {
        for (let z of H) (0, f9.checkReportMissingProp)(Q, z);
      });
      else X.if(S7._`${B} && (${(0, f9.checkMissingProp)(Q, H, G)})`), (0, f9.reportMissingProp)(Q, G), X.else();
    }
  }
  vz.validatePropertyDeps = _z;
  function kz(Q, $ = Q.schema) {
    let { gen: X, data: Y, keyword: J, it: G } = Q, W = X.name("valid");
    for (let H in $) {
      if ((0, xI.alwaysValidSchema)(G, $[H])) continue;
      X.if((0, f9.propertyInData)(X, Y, H, G.opts.ownProperties), () => {
        let B = Q.subschema({ keyword: J, schemaProp: H }, W);
        Q.mergeValidEvaluated(B, W);
      }, () => X.var(W, true)), Q.ok(W);
    }
  }
  vz.validateSchemaDeps = kz;
  vz.default = yI;
});
var fz = P((hz) => {
  Object.defineProperty(hz, "__esModule", { value: true });
  var gz = c(), uI = a(), mI = { message: "property name must be valid", params: ({ params: Q }) => gz._`{propertyName: ${Q.propertyName}}` }, lI = { keyword: "propertyNames", type: "object", schemaType: ["object", "boolean"], error: mI, code(Q) {
    let { gen: $, schema: X, data: Y, it: J } = Q;
    if ((0, uI.alwaysValidSchema)(J, X)) return;
    let G = $.name("valid");
    $.forIn("key", Y, (W) => {
      Q.setParams({ propertyName: W }), Q.subschema({ keyword: "propertyNames", data: W, dataTypes: ["string"], propertyName: W, compositeRule: true }, G), $.if((0, gz.not)(G), () => {
        if (Q.error(true), !J.allErrors) $.break();
      });
    }), Q.ok(G);
  } };
  hz.default = lI;
});
var _7 = P((uz) => {
  Object.defineProperty(uz, "__esModule", { value: true });
  var SQ = e0(), U1 = c(), pI = k1(), _Q = a(), dI = { message: "must NOT have additional properties", params: ({ params: Q }) => U1._`{additionalProperty: ${Q.additionalProperty}}` }, iI = { keyword: "additionalProperties", type: ["object"], schemaType: ["boolean", "object"], allowUndefined: true, trackErrors: true, error: dI, code(Q) {
    let { gen: $, schema: X, parentSchema: Y, data: J, errsCount: G, it: W } = Q;
    if (!G) throw Error("ajv implementation error");
    let { allErrors: H, opts: B } = W;
    if (W.props = true, B.removeAdditional !== "all" && (0, _Q.alwaysValidSchema)(W, X)) return;
    let z = (0, SQ.allSchemaProperties)(Y.properties), K = (0, SQ.allSchemaProperties)(Y.patternProperties);
    U(), Q.ok(U1._`${G} === ${pI.default.errors}`);
    function U() {
      $.forIn("key", J, (w) => {
        if (!z.length && !K.length) L(w);
        else $.if(q(w), () => L(w));
      });
    }
    function q(w) {
      let D;
      if (z.length > 8) {
        let M = (0, _Q.schemaRefOrVal)(W, Y.properties, "properties");
        D = (0, SQ.isOwnProperty)($, M, w);
      } else if (z.length) D = (0, U1.or)(...z.map((M) => U1._`${w} === ${M}`));
      else D = U1.nil;
      if (K.length) D = (0, U1.or)(D, ...K.map((M) => U1._`${(0, SQ.usePattern)(Q, M)}.test(${w})`));
      return (0, U1.not)(D);
    }
    function V(w) {
      $.code(U1._`delete ${J}[${w}]`);
    }
    function L(w) {
      if (B.removeAdditional === "all" || B.removeAdditional && X === false) {
        V(w);
        return;
      }
      if (X === false) {
        if (Q.setParams({ additionalProperty: w }), Q.error(), !H) $.break();
        return;
      }
      if (typeof X == "object" && !(0, _Q.alwaysValidSchema)(W, X)) {
        let D = $.name("valid");
        if (B.removeAdditional === "failing") F(w, D, false), $.if((0, U1.not)(D), () => {
          Q.reset(), V(w);
        });
        else if (F(w, D), !H) $.if((0, U1.not)(D), () => $.break());
      }
    }
    function F(w, D, M) {
      let R = { keyword: "additionalProperties", dataProp: w, dataPropType: _Q.Type.Str };
      if (M === false) Object.assign(R, { compositeRule: true, createErrors: false, allErrors: false });
      Q.subschema(R, D);
    }
  } };
  uz.default = iI;
});
var pz = P((cz) => {
  Object.defineProperty(cz, "__esModule", { value: true });
  var oI = b9(), mz = e0(), k7 = a(), lz = _7(), rI = { keyword: "properties", type: "object", schemaType: "object", code(Q) {
    let { gen: $, schema: X, parentSchema: Y, data: J, it: G } = Q;
    if (G.opts.removeAdditional === "all" && Y.additionalProperties === void 0) lz.default.code(new oI.KeywordCxt(G, lz.default, "additionalProperties"));
    let W = (0, mz.allSchemaProperties)(X);
    for (let U of W) G.definedProperties.add(U);
    if (G.opts.unevaluated && W.length && G.props !== true) G.props = k7.mergeEvaluated.props($, (0, k7.toHash)(W), G.props);
    let H = W.filter((U) => !(0, k7.alwaysValidSchema)(G, X[U]));
    if (H.length === 0) return;
    let B = $.name("valid");
    for (let U of H) {
      if (z(U)) K(U);
      else {
        if ($.if((0, mz.propertyInData)($, J, U, G.opts.ownProperties)), K(U), !G.allErrors) $.else().var(B, true);
        $.endIf();
      }
      Q.it.definedProperties.add(U), Q.ok(B);
    }
    function z(U) {
      return G.opts.useDefaults && !G.compositeRule && X[U].default !== void 0;
    }
    function K(U) {
      Q.subschema({ keyword: "properties", schemaProp: U, dataProp: U }, B);
    }
  } };
  cz.default = rI;
});
var rz = P((oz) => {
  Object.defineProperty(oz, "__esModule", { value: true });
  var dz = e0(), kQ = c(), iz = a(), nz = a(), aI = { keyword: "patternProperties", type: "object", schemaType: "object", code(Q) {
    let { gen: $, schema: X, data: Y, parentSchema: J, it: G } = Q, { opts: W } = G, H = (0, dz.allSchemaProperties)(X), B = H.filter((F) => (0, iz.alwaysValidSchema)(G, X[F]));
    if (H.length === 0 || B.length === H.length && (!G.opts.unevaluated || G.props === true)) return;
    let z = W.strictSchema && !W.allowMatchingProperties && J.properties, K = $.name("valid");
    if (G.props !== true && !(G.props instanceof kQ.Name)) G.props = (0, nz.evaluatedPropsToName)($, G.props);
    let { props: U } = G;
    q();
    function q() {
      for (let F of H) {
        if (z) V(F);
        if (G.allErrors) L(F);
        else $.var(K, true), L(F), $.if(K);
      }
    }
    function V(F) {
      for (let w in z) if (new RegExp(F).test(w)) (0, iz.checkStrictMode)(G, `property ${w} matches pattern ${F} (use allowMatchingProperties)`);
    }
    function L(F) {
      $.forIn("key", Y, (w) => {
        $.if(kQ._`${(0, dz.usePattern)(Q, F)}.test(${w})`, () => {
          let D = B.includes(F);
          if (!D) Q.subschema({ keyword: "patternProperties", schemaProp: F, dataProp: w, dataPropType: nz.Type.Str }, K);
          if (G.opts.unevaluated && U !== true) $.assign(kQ._`${U}[${w}]`, true);
          else if (!D && !G.allErrors) $.if((0, kQ.not)(K), () => $.break());
        });
      });
    }
  } };
  oz.default = aI;
});
var az = P((tz) => {
  Object.defineProperty(tz, "__esModule", { value: true });
  var eI = a(), QE = { keyword: "not", schemaType: ["object", "boolean"], trackErrors: true, code(Q) {
    let { gen: $, schema: X, it: Y } = Q;
    if ((0, eI.alwaysValidSchema)(Y, X)) {
      Q.fail();
      return;
    }
    let J = $.name("valid");
    Q.subschema({ keyword: "not", compositeRule: true, createErrors: false, allErrors: false }, J), Q.failResult(J, () => Q.reset(), () => Q.error());
  }, error: { message: "must NOT be valid" } };
  tz.default = QE;
});
var ez = P((sz) => {
  Object.defineProperty(sz, "__esModule", { value: true });
  var XE = e0(), YE = { keyword: "anyOf", schemaType: "array", trackErrors: true, code: XE.validateUnion, error: { message: "must match a schema in anyOf" } };
  sz.default = YE;
});
var $K = P((QK) => {
  Object.defineProperty(QK, "__esModule", { value: true });
  var vQ = c(), GE = a(), WE = { message: "must match exactly one schema in oneOf", params: ({ params: Q }) => vQ._`{passingSchemas: ${Q.passing}}` }, HE = { keyword: "oneOf", schemaType: "array", trackErrors: true, error: WE, code(Q) {
    let { gen: $, schema: X, parentSchema: Y, it: J } = Q;
    if (!Array.isArray(X)) throw Error("ajv implementation error");
    if (J.opts.discriminator && Y.discriminator) return;
    let G = X, W = $.let("valid", false), H = $.let("passing", null), B = $.name("_valid");
    Q.setParams({ passing: H }), $.block(z), Q.result(W, () => Q.reset(), () => Q.error(true));
    function z() {
      G.forEach((K, U) => {
        let q;
        if ((0, GE.alwaysValidSchema)(J, K)) $.var(B, true);
        else q = Q.subschema({ keyword: "oneOf", schemaProp: U, compositeRule: true }, B);
        if (U > 0) $.if(vQ._`${B} && ${W}`).assign(W, false).assign(H, vQ._`[${H}, ${U}]`).else();
        $.if(B, () => {
          if ($.assign(W, true), $.assign(H, U), q) Q.mergeEvaluated(q, vQ.Name);
        });
      });
    }
  } };
  QK.default = HE;
});
var YK = P((XK) => {
  Object.defineProperty(XK, "__esModule", { value: true });
  var zE = a(), KE = { keyword: "allOf", schemaType: "array", code(Q) {
    let { gen: $, schema: X, it: Y } = Q;
    if (!Array.isArray(X)) throw Error("ajv implementation error");
    let J = $.name("valid");
    X.forEach((G, W) => {
      if ((0, zE.alwaysValidSchema)(Y, G)) return;
      let H = Q.subschema({ keyword: "allOf", schemaProp: W }, J);
      Q.ok(J), Q.mergeEvaluated(H);
    });
  } };
  XK.default = KE;
});
var HK = P((WK) => {
  Object.defineProperty(WK, "__esModule", { value: true });
  var TQ = c(), GK = a(), qE = { message: ({ params: Q }) => TQ.str`must match "${Q.ifClause}" schema`, params: ({ params: Q }) => TQ._`{failingKeyword: ${Q.ifClause}}` }, UE = { keyword: "if", schemaType: ["object", "boolean"], trackErrors: true, error: qE, code(Q) {
    let { gen: $, parentSchema: X, it: Y } = Q;
    if (X.then === void 0 && X.else === void 0) (0, GK.checkStrictMode)(Y, '"if" without "then" and "else" is ignored');
    let J = JK(Y, "then"), G = JK(Y, "else");
    if (!J && !G) return;
    let W = $.let("valid", true), H = $.name("_valid");
    if (B(), Q.reset(), J && G) {
      let K = $.let("ifClause");
      Q.setParams({ ifClause: K }), $.if(H, z("then", K), z("else", K));
    } else if (J) $.if(H, z("then"));
    else $.if((0, TQ.not)(H), z("else"));
    Q.pass(W, () => Q.error(true));
    function B() {
      let K = Q.subschema({ keyword: "if", compositeRule: true, createErrors: false, allErrors: false }, H);
      Q.mergeEvaluated(K);
    }
    function z(K, U) {
      return () => {
        let q = Q.subschema({ keyword: K }, H);
        if ($.assign(W, H), Q.mergeValidEvaluated(q, W), U) $.assign(U, TQ._`${K}`);
        else Q.setParams({ ifClause: K });
      };
    }
  } };
  function JK(Q, $) {
    let X = Q.schema[$];
    return X !== void 0 && !(0, GK.alwaysValidSchema)(Q, X);
  }
  WK.default = UE;
});
var zK = P((BK) => {
  Object.defineProperty(BK, "__esModule", { value: true });
  var FE = a(), NE = { keyword: ["then", "else"], schemaType: ["object", "boolean"], code({ keyword: Q, parentSchema: $, it: X }) {
    if ($.if === void 0) (0, FE.checkStrictMode)(X, `"${Q}" without "if" is ignored`);
  } };
  BK.default = NE;
});
var VK = P((KK) => {
  Object.defineProperty(KK, "__esModule", { value: true });
  var DE = Z7(), wE = Ez(), ME = C7(), AE = Zz(), jE = Sz(), RE = yz(), IE = fz(), EE = _7(), PE = pz(), bE = rz(), ZE = az(), CE = ez(), SE = $K(), _E = YK(), kE = HK(), vE = zK();
  function TE(Q = false) {
    let $ = [ZE.default, CE.default, SE.default, _E.default, kE.default, vE.default, IE.default, EE.default, RE.default, PE.default, bE.default];
    if (Q) $.push(wE.default, AE.default);
    else $.push(DE.default, ME.default);
    return $.push(jE.default), $;
  }
  KK.default = TE;
});
var UK = P((qK) => {
  Object.defineProperty(qK, "__esModule", { value: true });
  var N0 = c(), yE = { message: ({ schemaCode: Q }) => N0.str`must match format "${Q}"`, params: ({ schemaCode: Q }) => N0._`{format: ${Q}}` }, gE = { keyword: "format", type: ["number", "string"], schemaType: "string", $data: true, error: yE, code(Q, $) {
    let { gen: X, data: Y, $data: J, schema: G, schemaCode: W, it: H } = Q, { opts: B, errSchemaPath: z, schemaEnv: K, self: U } = H;
    if (!B.validateFormats) return;
    if (J) q();
    else V();
    function q() {
      let L = X.scopeValue("formats", { ref: U.formats, code: B.code.formats }), F = X.const("fDef", N0._`${L}[${W}]`), w = X.let("fType"), D = X.let("format");
      X.if(N0._`typeof ${F} == "object" && !(${F} instanceof RegExp)`, () => X.assign(w, N0._`${F}.type || "string"`).assign(D, N0._`${F}.validate`), () => X.assign(w, N0._`"string"`).assign(D, F)), Q.fail$data((0, N0.or)(M(), R()));
      function M() {
        if (B.strictSchema === false) return N0.nil;
        return N0._`${W} && !${D}`;
      }
      function R() {
        let Z = K.$async ? N0._`(${F}.async ? await ${D}(${Y}) : ${D}(${Y}))` : N0._`${D}(${Y})`, v = N0._`(typeof ${D} == "function" ? ${Z} : ${D}.test(${Y}))`;
        return N0._`${D} && ${D} !== true && ${w} === ${$} && !${v}`;
      }
    }
    function V() {
      let L = U.formats[G];
      if (!L) {
        M();
        return;
      }
      if (L === true) return;
      let [F, w, D] = R(L);
      if (F === $) Q.pass(Z());
      function M() {
        if (B.strictSchema === false) {
          U.logger.warn(v());
          return;
        }
        throw Error(v());
        function v() {
          return `unknown format "${G}" ignored in schema at path "${z}"`;
        }
      }
      function R(v) {
        let O0 = v instanceof RegExp ? (0, N0.regexpCode)(v) : B.code.formats ? N0._`${B.code.formats}${(0, N0.getProperty)(G)}` : void 0, D0 = X.scopeValue("formats", { key: G, ref: v, code: O0 });
        if (typeof v == "object" && !(v instanceof RegExp)) return [v.type || "string", v.validate, N0._`${D0}.validate`];
        return ["string", v, D0];
      }
      function Z() {
        if (typeof L == "object" && !(L instanceof RegExp) && L.async) {
          if (!K.$async) throw Error("async format in sync schema");
          return N0._`await ${D}(${Y})`;
        }
        return typeof w == "function" ? N0._`${D}(${Y})` : N0._`${D}.test(${Y})`;
      }
    }
  } };
  qK.default = gE;
});
var FK = P((LK) => {
  Object.defineProperty(LK, "__esModule", { value: true });
  var fE = UK(), uE = [fE.default];
  LK.default = uE;
});
var DK = P((NK) => {
  Object.defineProperty(NK, "__esModule", { value: true });
  NK.contentVocabulary = NK.metadataVocabulary = void 0;
  NK.metadataVocabulary = ["title", "description", "default", "deprecated", "readOnly", "writeOnly", "examples"];
  NK.contentVocabulary = ["contentMediaType", "contentEncoding", "contentSchema"];
});
var AK = P((MK) => {
  Object.defineProperty(MK, "__esModule", { value: true });
  var cE = cB(), pE = Nz(), dE = VK(), iE = FK(), wK = DK(), nE = [cE.default, pE.default, (0, dE.default)(), iE.default, wK.metadataVocabulary, wK.contentVocabulary];
  MK.default = nE;
});
var EK = P((RK) => {
  Object.defineProperty(RK, "__esModule", { value: true });
  RK.DiscrError = void 0;
  var jK;
  (function(Q) {
    Q.Tag = "tag", Q.Mapping = "mapping";
  })(jK || (RK.DiscrError = jK = {}));
});
var ZK = P((bK) => {
  Object.defineProperty(bK, "__esModule", { value: true });
  var b4 = c(), v7 = EK(), PK = NQ(), rE = Z9(), tE = a(), aE = { message: ({ params: { discrError: Q, tagName: $ } }) => Q === v7.DiscrError.Tag ? `tag "${$}" must be string` : `value of tag "${$}" must be in oneOf`, params: ({ params: { discrError: Q, tag: $, tagName: X } }) => b4._`{error: ${Q}, tag: ${X}, tagValue: ${$}}` }, sE = { keyword: "discriminator", type: "object", schemaType: "object", error: aE, code(Q) {
    let { gen: $, data: X, schema: Y, parentSchema: J, it: G } = Q, { oneOf: W } = J;
    if (!G.opts.discriminator) throw Error("discriminator: requires discriminator option");
    let H = Y.propertyName;
    if (typeof H != "string") throw Error("discriminator: requires propertyName");
    if (Y.mapping) throw Error("discriminator: mapping is not supported");
    if (!W) throw Error("discriminator: requires oneOf keyword");
    let B = $.let("valid", false), z = $.const("tag", b4._`${X}${(0, b4.getProperty)(H)}`);
    $.if(b4._`typeof ${z} == "string"`, () => K(), () => Q.error(false, { discrError: v7.DiscrError.Tag, tag: z, tagName: H })), Q.ok(B);
    function K() {
      let V = q();
      $.if(false);
      for (let L in V) $.elseIf(b4._`${z} === ${L}`), $.assign(B, U(V[L]));
      $.else(), Q.error(false, { discrError: v7.DiscrError.Mapping, tag: z, tagName: H }), $.endIf();
    }
    function U(V) {
      let L = $.name("valid"), F = Q.subschema({ keyword: "oneOf", schemaProp: V }, L);
      return Q.mergeEvaluated(F, b4.Name), L;
    }
    function q() {
      var V;
      let L = {}, F = D(J), w = true;
      for (let Z = 0; Z < W.length; Z++) {
        let v = W[Z];
        if ((v === null || v === void 0 ? void 0 : v.$ref) && !(0, tE.schemaHasRulesButRef)(v, G.self.RULES)) {
          let D0 = v.$ref;
          if (v = PK.resolveRef.call(G.self, G.schemaEnv.root, G.baseId, D0), v instanceof PK.SchemaEnv) v = v.schema;
          if (v === void 0) throw new rE.default(G.opts.uriResolver, G.baseId, D0);
        }
        let O0 = (V = v === null || v === void 0 ? void 0 : v.properties) === null || V === void 0 ? void 0 : V[H];
        if (typeof O0 != "object") throw Error(`discriminator: oneOf subschemas (or referenced schemas) must have "properties/${H}"`);
        w = w && (F || D(v)), M(O0, Z);
      }
      if (!w) throw Error(`discriminator: "${H}" must be required`);
      return L;
      function D({ required: Z }) {
        return Array.isArray(Z) && Z.includes(H);
      }
      function M(Z, v) {
        if (Z.const) R(Z.const, v);
        else if (Z.enum) for (let O0 of Z.enum) R(O0, v);
        else throw Error(`discriminator: "properties/${H}" must have "const" or "enum"`);
      }
      function R(Z, v) {
        if (typeof Z != "string" || Z in L) throw Error(`discriminator: "${H}" values must be unique strings`);
        L[Z] = v;
      }
    }
  } };
  bK.default = sE;
});
var CK = P((wg, QP) => {
  QP.exports = { $schema: "http://json-schema.org/draft-07/schema#", $id: "http://json-schema.org/draft-07/schema#", title: "Core schema meta-schema", definitions: { schemaArray: { type: "array", minItems: 1, items: { $ref: "#" } }, nonNegativeInteger: { type: "integer", minimum: 0 }, nonNegativeIntegerDefault0: { allOf: [{ $ref: "#/definitions/nonNegativeInteger" }, { default: 0 }] }, simpleTypes: { enum: ["array", "boolean", "integer", "null", "number", "object", "string"] }, stringArray: { type: "array", items: { type: "string" }, uniqueItems: true, default: [] } }, type: ["object", "boolean"], properties: { $id: { type: "string", format: "uri-reference" }, $schema: { type: "string", format: "uri" }, $ref: { type: "string", format: "uri-reference" }, $comment: { type: "string" }, title: { type: "string" }, description: { type: "string" }, default: true, readOnly: { type: "boolean", default: false }, examples: { type: "array", items: true }, multipleOf: { type: "number", exclusiveMinimum: 0 }, maximum: { type: "number" }, exclusiveMaximum: { type: "number" }, minimum: { type: "number" }, exclusiveMinimum: { type: "number" }, maxLength: { $ref: "#/definitions/nonNegativeInteger" }, minLength: { $ref: "#/definitions/nonNegativeIntegerDefault0" }, pattern: { type: "string", format: "regex" }, additionalItems: { $ref: "#" }, items: { anyOf: [{ $ref: "#" }, { $ref: "#/definitions/schemaArray" }], default: true }, maxItems: { $ref: "#/definitions/nonNegativeInteger" }, minItems: { $ref: "#/definitions/nonNegativeIntegerDefault0" }, uniqueItems: { type: "boolean", default: false }, contains: { $ref: "#" }, maxProperties: { $ref: "#/definitions/nonNegativeInteger" }, minProperties: { $ref: "#/definitions/nonNegativeIntegerDefault0" }, required: { $ref: "#/definitions/stringArray" }, additionalProperties: { $ref: "#" }, definitions: { type: "object", additionalProperties: { $ref: "#" }, default: {} }, properties: { type: "object", additionalProperties: { $ref: "#" }, default: {} }, patternProperties: { type: "object", additionalProperties: { $ref: "#" }, propertyNames: { format: "regex" }, default: {} }, dependencies: { type: "object", additionalProperties: { anyOf: [{ $ref: "#" }, { $ref: "#/definitions/stringArray" }] } }, propertyNames: { $ref: "#" }, const: true, enum: { type: "array", items: true, minItems: 1, uniqueItems: true }, type: { anyOf: [{ $ref: "#/definitions/simpleTypes" }, { type: "array", items: { $ref: "#/definitions/simpleTypes" }, minItems: 1, uniqueItems: true }] }, format: { type: "string" }, contentMediaType: { type: "string" }, contentEncoding: { type: "string" }, if: { $ref: "#" }, then: { $ref: "#" }, else: { $ref: "#" }, allOf: { $ref: "#/definitions/schemaArray" }, anyOf: { $ref: "#/definitions/schemaArray" }, oneOf: { $ref: "#/definitions/schemaArray" }, not: { $ref: "#" } }, default: true };
});
var x7 = P((u0, T7) => {
  Object.defineProperty(u0, "__esModule", { value: true });
  u0.MissingRefError = u0.ValidationError = u0.CodeGen = u0.Name = u0.nil = u0.stringify = u0.str = u0._ = u0.KeywordCxt = u0.Ajv = void 0;
  var $P = vB(), XP = AK(), YP = ZK(), SK = CK(), JP = ["/properties"], xQ = "http://json-schema.org/draft-07/schema";
  class u9 extends $P.default {
    _addVocabularies() {
      if (super._addVocabularies(), XP.default.forEach((Q) => this.addVocabulary(Q)), this.opts.discriminator) this.addKeyword(YP.default);
    }
    _addDefaultMetaSchema() {
      if (super._addDefaultMetaSchema(), !this.opts.meta) return;
      let Q = this.opts.$data ? this.$dataMetaSchema(SK, JP) : SK;
      this.addMetaSchema(Q, xQ, false), this.refs["http://json-schema.org/schema"] = xQ;
    }
    defaultMeta() {
      return this.opts.defaultMeta = super.defaultMeta() || (this.getSchema(xQ) ? xQ : void 0);
    }
  }
  u0.Ajv = u9;
  T7.exports = u0 = u9;
  T7.exports.Ajv = u9;
  Object.defineProperty(u0, "__esModule", { value: true });
  u0.default = u9;
  var GP = b9();
  Object.defineProperty(u0, "KeywordCxt", { enumerable: true, get: function() {
    return GP.KeywordCxt;
  } });
  var Z4 = c();
  Object.defineProperty(u0, "_", { enumerable: true, get: function() {
    return Z4._;
  } });
  Object.defineProperty(u0, "str", { enumerable: true, get: function() {
    return Z4.str;
  } });
  Object.defineProperty(u0, "stringify", { enumerable: true, get: function() {
    return Z4.stringify;
  } });
  Object.defineProperty(u0, "nil", { enumerable: true, get: function() {
    return Z4.nil;
  } });
  Object.defineProperty(u0, "Name", { enumerable: true, get: function() {
    return Z4.Name;
  } });
  Object.defineProperty(u0, "CodeGen", { enumerable: true, get: function() {
    return Z4.CodeGen;
  } });
  var WP = LQ();
  Object.defineProperty(u0, "ValidationError", { enumerable: true, get: function() {
    return WP.default;
  } });
  var HP = Z9();
  Object.defineProperty(u0, "MissingRefError", { enumerable: true, get: function() {
    return HP.default;
  } });
});
var uK = P((hK) => {
  Object.defineProperty(hK, "__esModule", { value: true });
  hK.formatNames = hK.fastFormats = hK.fullFormats = void 0;
  function E1(Q, $) {
    return { validate: Q, compare: $ };
  }
  hK.fullFormats = { date: E1(TK, f7), time: E1(g7(true), u7), "date-time": E1(_K(true), yK), "iso-time": E1(g7(), xK), "iso-date-time": E1(_K(), gK), duration: /^P(?!$)((\d+Y)?(\d+M)?(\d+D)?(T(?=\d)(\d+H)?(\d+M)?(\d+S)?)?|(\d+W)?)$/, uri: FP, "uri-reference": /^(?:[a-z][a-z0-9+\-.]*:)?(?:\/?\/(?:(?:[a-z0-9\-._~!$&'()*+,;=:]|%[0-9a-f]{2})*@)?(?:\[(?:(?:(?:(?:[0-9a-f]{1,4}:){6}|::(?:[0-9a-f]{1,4}:){5}|(?:[0-9a-f]{1,4})?::(?:[0-9a-f]{1,4}:){4}|(?:(?:[0-9a-f]{1,4}:){0,1}[0-9a-f]{1,4})?::(?:[0-9a-f]{1,4}:){3}|(?:(?:[0-9a-f]{1,4}:){0,2}[0-9a-f]{1,4})?::(?:[0-9a-f]{1,4}:){2}|(?:(?:[0-9a-f]{1,4}:){0,3}[0-9a-f]{1,4})?::[0-9a-f]{1,4}:|(?:(?:[0-9a-f]{1,4}:){0,4}[0-9a-f]{1,4})?::)(?:[0-9a-f]{1,4}:[0-9a-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?))|(?:(?:[0-9a-f]{1,4}:){0,5}[0-9a-f]{1,4})?::[0-9a-f]{1,4}|(?:(?:[0-9a-f]{1,4}:){0,6}[0-9a-f]{1,4})?::)|[Vv][0-9a-f]+\.[a-z0-9\-._~!$&'()*+,;=:]+)\]|(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?)|(?:[a-z0-9\-._~!$&'"()*+,;=]|%[0-9a-f]{2})*)(?::\d*)?(?:\/(?:[a-z0-9\-._~!$&'"()*+,;=:@]|%[0-9a-f]{2})*)*|\/(?:(?:[a-z0-9\-._~!$&'"()*+,;=:@]|%[0-9a-f]{2})+(?:\/(?:[a-z0-9\-._~!$&'"()*+,;=:@]|%[0-9a-f]{2})*)*)?|(?:[a-z0-9\-._~!$&'"()*+,;=:@]|%[0-9a-f]{2})+(?:\/(?:[a-z0-9\-._~!$&'"()*+,;=:@]|%[0-9a-f]{2})*)*)?(?:\?(?:[a-z0-9\-._~!$&'"()*+,;=:@/?]|%[0-9a-f]{2})*)?(?:#(?:[a-z0-9\-._~!$&'"()*+,;=:@/?]|%[0-9a-f]{2})*)?$/i, "uri-template": /^(?:(?:[^\x00-\x20"'<>%\\^`{|}]|%[0-9a-f]{2})|\{[+#./;?&=,!@|]?(?:[a-z0-9_]|%[0-9a-f]{2})+(?::[1-9][0-9]{0,3}|\*)?(?:,(?:[a-z0-9_]|%[0-9a-f]{2})+(?::[1-9][0-9]{0,3}|\*)?)*\})*$/i, url: /^(?:https?|ftp):\/\/(?:\S+(?::\S*)?@)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z0-9\u{00a1}-\u{ffff}]+-)*[a-z0-9\u{00a1}-\u{ffff}]+)(?:\.(?:[a-z0-9\u{00a1}-\u{ffff}]+-)*[a-z0-9\u{00a1}-\u{ffff}]+)*(?:\.(?:[a-z\u{00a1}-\u{ffff}]{2,})))(?::\d{2,5})?(?:\/[^\s]*)?$/iu, email: /^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/i, hostname: /^(?=.{1,253}\.?$)[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[-0-9a-z]{0,61}[0-9a-z])?)*\.?$/i, ipv4: /^(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)$/, ipv6: /^((([0-9a-f]{1,4}:){7}([0-9a-f]{1,4}|:))|(([0-9a-f]{1,4}:){6}(:[0-9a-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9a-f]{1,4}:){5}(((:[0-9a-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9a-f]{1,4}:){4}(((:[0-9a-f]{1,4}){1,3})|((:[0-9a-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9a-f]{1,4}:){3}(((:[0-9a-f]{1,4}){1,4})|((:[0-9a-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9a-f]{1,4}:){2}(((:[0-9a-f]{1,4}){1,5})|((:[0-9a-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9a-f]{1,4}:){1}(((:[0-9a-f]{1,4}){1,6})|((:[0-9a-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9a-f]{1,4}){1,7})|((:[0-9a-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))$/i, regex: jP, uuid: /^(?:urn:uuid:)?[0-9a-f]{8}-(?:[0-9a-f]{4}-){3}[0-9a-f]{12}$/i, "json-pointer": /^(?:\/(?:[^~/]|~0|~1)*)*$/, "json-pointer-uri-fragment": /^#(?:\/(?:[a-z0-9_\-.!$&'()*+,;:=@]|%[0-9a-f]{2}|~0|~1)*)*$/i, "relative-json-pointer": /^(?:0|[1-9][0-9]*)(?:#|(?:\/(?:[^~/]|~0|~1)*)*)$/, byte: NP, int32: { type: "number", validate: wP }, int64: { type: "number", validate: MP }, float: { type: "number", validate: vK }, double: { type: "number", validate: vK }, password: true, binary: true };
  hK.fastFormats = { ...hK.fullFormats, date: E1(/^\d\d\d\d-[0-1]\d-[0-3]\d$/, f7), time: E1(/^(?:[0-2]\d:[0-5]\d:[0-5]\d|23:59:60)(?:\.\d+)?(?:z|[+-]\d\d(?::?\d\d)?)$/i, u7), "date-time": E1(/^\d\d\d\d-[0-1]\d-[0-3]\dt(?:[0-2]\d:[0-5]\d:[0-5]\d|23:59:60)(?:\.\d+)?(?:z|[+-]\d\d(?::?\d\d)?)$/i, yK), "iso-time": E1(/^(?:[0-2]\d:[0-5]\d:[0-5]\d|23:59:60)(?:\.\d+)?(?:z|[+-]\d\d(?::?\d\d)?)?$/i, xK), "iso-date-time": E1(/^\d\d\d\d-[0-1]\d-[0-3]\d[t\s](?:[0-2]\d:[0-5]\d:[0-5]\d|23:59:60)(?:\.\d+)?(?:z|[+-]\d\d(?::?\d\d)?)?$/i, gK), uri: /^(?:[a-z][a-z0-9+\-.]*:)(?:\/?\/)?[^\s]*$/i, "uri-reference": /^(?:(?:[a-z][a-z0-9+\-.]*:)?\/?\/)?(?:[^\\\s#][^\s#]*)?(?:#[^\\\s]*)?$/i, email: /^[a-z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)*$/i };
  hK.formatNames = Object.keys(hK.fullFormats);
  function KP(Q) {
    return Q % 4 === 0 && (Q % 100 !== 0 || Q % 400 === 0);
  }
  var VP = /^(\d\d\d\d)-(\d\d)-(\d\d)$/, qP = [0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];
  function TK(Q) {
    let $ = VP.exec(Q);
    if (!$) return false;
    let X = +$[1], Y = +$[2], J = +$[3];
    return Y >= 1 && Y <= 12 && J >= 1 && J <= (Y === 2 && KP(X) ? 29 : qP[Y]);
  }
  function f7(Q, $) {
    if (!(Q && $)) return;
    if (Q > $) return 1;
    if (Q < $) return -1;
    return 0;
  }
  var y7 = /^(\d\d):(\d\d):(\d\d(?:\.\d+)?)(z|([+-])(\d\d)(?::?(\d\d))?)?$/i;
  function g7(Q) {
    return function(X) {
      let Y = y7.exec(X);
      if (!Y) return false;
      let J = +Y[1], G = +Y[2], W = +Y[3], H = Y[4], B = Y[5] === "-" ? -1 : 1, z = +(Y[6] || 0), K = +(Y[7] || 0);
      if (z > 23 || K > 59 || Q && !H) return false;
      if (J <= 23 && G <= 59 && W < 60) return true;
      let U = G - K * B, q = J - z * B - (U < 0 ? 1 : 0);
      return (q === 23 || q === -1) && (U === 59 || U === -1) && W < 61;
    };
  }
  function u7(Q, $) {
    if (!(Q && $)) return;
    let X = (/* @__PURE__ */ new Date("2020-01-01T" + Q)).valueOf(), Y = (/* @__PURE__ */ new Date("2020-01-01T" + $)).valueOf();
    if (!(X && Y)) return;
    return X - Y;
  }
  function xK(Q, $) {
    if (!(Q && $)) return;
    let X = y7.exec(Q), Y = y7.exec($);
    if (!(X && Y)) return;
    if (Q = X[1] + X[2] + X[3], $ = Y[1] + Y[2] + Y[3], Q > $) return 1;
    if (Q < $) return -1;
    return 0;
  }
  var h7 = /t|\s/i;
  function _K(Q) {
    let $ = g7(Q);
    return function(Y) {
      let J = Y.split(h7);
      return J.length === 2 && TK(J[0]) && $(J[1]);
    };
  }
  function yK(Q, $) {
    if (!(Q && $)) return;
    let X = new Date(Q).valueOf(), Y = new Date($).valueOf();
    if (!(X && Y)) return;
    return X - Y;
  }
  function gK(Q, $) {
    if (!(Q && $)) return;
    let [X, Y] = Q.split(h7), [J, G] = $.split(h7), W = f7(X, J);
    if (W === void 0) return;
    return W || u7(Y, G);
  }
  var UP = /\/|:/, LP = /^(?:[a-z][a-z0-9+\-.]*:)(?:\/?\/(?:(?:[a-z0-9\-._~!$&'()*+,;=:]|%[0-9a-f]{2})*@)?(?:\[(?:(?:(?:(?:[0-9a-f]{1,4}:){6}|::(?:[0-9a-f]{1,4}:){5}|(?:[0-9a-f]{1,4})?::(?:[0-9a-f]{1,4}:){4}|(?:(?:[0-9a-f]{1,4}:){0,1}[0-9a-f]{1,4})?::(?:[0-9a-f]{1,4}:){3}|(?:(?:[0-9a-f]{1,4}:){0,2}[0-9a-f]{1,4})?::(?:[0-9a-f]{1,4}:){2}|(?:(?:[0-9a-f]{1,4}:){0,3}[0-9a-f]{1,4})?::[0-9a-f]{1,4}:|(?:(?:[0-9a-f]{1,4}:){0,4}[0-9a-f]{1,4})?::)(?:[0-9a-f]{1,4}:[0-9a-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?))|(?:(?:[0-9a-f]{1,4}:){0,5}[0-9a-f]{1,4})?::[0-9a-f]{1,4}|(?:(?:[0-9a-f]{1,4}:){0,6}[0-9a-f]{1,4})?::)|[Vv][0-9a-f]+\.[a-z0-9\-._~!$&'()*+,;=:]+)\]|(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?)|(?:[a-z0-9\-._~!$&'()*+,;=]|%[0-9a-f]{2})*)(?::\d*)?(?:\/(?:[a-z0-9\-._~!$&'()*+,;=:@]|%[0-9a-f]{2})*)*|\/(?:(?:[a-z0-9\-._~!$&'()*+,;=:@]|%[0-9a-f]{2})+(?:\/(?:[a-z0-9\-._~!$&'()*+,;=:@]|%[0-9a-f]{2})*)*)?|(?:[a-z0-9\-._~!$&'()*+,;=:@]|%[0-9a-f]{2})+(?:\/(?:[a-z0-9\-._~!$&'()*+,;=:@]|%[0-9a-f]{2})*)*)(?:\?(?:[a-z0-9\-._~!$&'()*+,;=:@/?]|%[0-9a-f]{2})*)?(?:#(?:[a-z0-9\-._~!$&'()*+,;=:@/?]|%[0-9a-f]{2})*)?$/i;
  function FP(Q) {
    return UP.test(Q) && LP.test(Q);
  }
  var kK = /^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/gm;
  function NP(Q) {
    return kK.lastIndex = 0, kK.test(Q);
  }
  var OP = -2147483648, DP = 2147483647;
  function wP(Q) {
    return Number.isInteger(Q) && Q <= DP && Q >= OP;
  }
  function MP(Q) {
    return Number.isInteger(Q);
  }
  function vK() {
    return true;
  }
  var AP = /[^\\]\\Z/;
  function jP(Q) {
    if (AP.test(Q)) return false;
    try {
      return new RegExp(Q), true;
    } catch ($) {
      return false;
    }
  }
});
var lK = P((mK) => {
  Object.defineProperty(mK, "__esModule", { value: true });
  mK.formatLimitDefinition = void 0;
  var IP = x7(), L1 = c(), W6 = L1.operators, yQ = { formatMaximum: { okStr: "<=", ok: W6.LTE, fail: W6.GT }, formatMinimum: { okStr: ">=", ok: W6.GTE, fail: W6.LT }, formatExclusiveMaximum: { okStr: "<", ok: W6.LT, fail: W6.GTE }, formatExclusiveMinimum: { okStr: ">", ok: W6.GT, fail: W6.LTE } }, EP = { message: ({ keyword: Q, schemaCode: $ }) => L1.str`should be ${yQ[Q].okStr} ${$}`, params: ({ keyword: Q, schemaCode: $ }) => L1._`{comparison: ${yQ[Q].okStr}, limit: ${$}}` };
  mK.formatLimitDefinition = { keyword: Object.keys(yQ), type: "string", schemaType: "string", $data: true, error: EP, code(Q) {
    let { gen: $, data: X, schemaCode: Y, keyword: J, it: G } = Q, { opts: W, self: H } = G;
    if (!W.validateFormats) return;
    let B = new IP.KeywordCxt(G, H.RULES.all.format.definition, "format");
    if (B.$data) z();
    else K();
    function z() {
      let q = $.scopeValue("formats", { ref: H.formats, code: W.code.formats }), V = $.const("fmt", L1._`${q}[${B.schemaCode}]`);
      Q.fail$data((0, L1.or)(L1._`typeof ${V} != "object"`, L1._`${V} instanceof RegExp`, L1._`typeof ${V}.compare != "function"`, U(V)));
    }
    function K() {
      let q = B.schema, V = H.formats[q];
      if (!V || V === true) return;
      if (typeof V != "object" || V instanceof RegExp || typeof V.compare != "function") throw Error(`"${J}": format "${q}" does not define "compare" function`);
      let L = $.scopeValue("formats", { key: q, ref: V, code: W.code.formats ? L1._`${W.code.formats}${(0, L1.getProperty)(q)}` : void 0 });
      Q.fail$data(U(L));
    }
    function U(q) {
      return L1._`${q}.compare(${X}, ${Y}) ${yQ[J].fail} 0`;
    }
  }, dependencies: ["format"] };
  var PP = (Q) => {
    return Q.addKeyword(mK.formatLimitDefinition), Q;
  };
  mK.default = PP;
});
var iK = P((m9, dK) => {
  Object.defineProperty(m9, "__esModule", { value: true });
  var C4 = uK(), ZP = lK(), c7 = c(), cK = new c7.Name("fullFormats"), CP = new c7.Name("fastFormats"), p7 = (Q, $ = { keywords: true }) => {
    if (Array.isArray($)) return pK(Q, $, C4.fullFormats, cK), Q;
    let [X, Y] = $.mode === "fast" ? [C4.fastFormats, CP] : [C4.fullFormats, cK], J = $.formats || C4.formatNames;
    if (pK(Q, J, X, Y), $.keywords) (0, ZP.default)(Q);
    return Q;
  };
  p7.get = (Q, $ = "full") => {
    let Y = ($ === "fast" ? C4.fastFormats : C4.fullFormats)[Q];
    if (!Y) throw Error(`Unknown format "${Q}"`);
    return Y;
  };
  function pK(Q, $, X, Y) {
    var J, G;
    (J = (G = Q.opts.code).formats) !== null && J !== void 0 || (G.formats = c7._`require("ajv-formats/dist/formats").${Y}`);
    for (let W of $) Q.addFormat(W, X[W]);
  }
  dK.exports = m9 = p7;
  Object.defineProperty(m9, "__esModule", { value: true });
  m9.default = p7;
});
var kV = 50;
function g6(Q = kV) {
  let $ = new AbortController();
  return _V(Q, $.signal), $;
}
var T0 = class extends Error {
};
function h6() {
  return process.versions.bun !== void 0;
}
var xV = typeof global == "object" && global && global.Object === Object && global;
var W5 = xV;
var yV = typeof self == "object" && self && self.Object === Object && self;
var gV = W5 || yV || Function("return this")();
var f6 = gV;
var hV = f6.Symbol;
var u6 = hV;
var H5 = Object.prototype;
var fV = H5.hasOwnProperty;
var uV = H5.toString;
var T4 = u6 ? u6.toStringTag : void 0;
function mV(Q) {
  var $ = fV.call(Q, T4), X = Q[T4];
  try {
    Q[T4] = void 0;
    var Y = true;
  } catch (G) {
  }
  var J = uV.call(Q);
  if (Y) if ($) Q[T4] = X;
  else delete Q[T4];
  return J;
}
var B5 = mV;
var lV = Object.prototype;
var cV = lV.toString;
function pV(Q) {
  return cV.call(Q);
}
var z5 = pV;
var dV = "[object Null]";
var iV = "[object Undefined]";
var K5 = u6 ? u6.toStringTag : void 0;
function nV(Q) {
  if (Q == null) return Q === void 0 ? iV : dV;
  return K5 && K5 in Object(Q) ? B5(Q) : z5(Q);
}
var V5 = nV;
function oV(Q) {
  var $ = typeof Q;
  return Q != null && ($ == "object" || $ == "function");
}
var p9 = oV;
var rV = "[object AsyncFunction]";
var tV = "[object Function]";
var aV = "[object GeneratorFunction]";
var sV = "[object Proxy]";
function eV(Q) {
  if (!p9(Q)) return false;
  var $ = V5(Q);
  return $ == tV || $ == aV || $ == rV || $ == sV;
}
var q5 = eV;
var Qq = f6["__core-js_shared__"];
var d9 = Qq;
var U5 = (function() {
  var Q = /[^.]+$/.exec(d9 && d9.keys && d9.keys.IE_PROTO || "");
  return Q ? "Symbol(src)_1." + Q : "";
})();
function $q(Q) {
  return !!U5 && U5 in Q;
}
var L5 = $q;
var Xq = Function.prototype;
var Yq = Xq.toString;
function Jq(Q) {
  if (Q != null) {
    try {
      return Yq.call(Q);
    } catch ($) {
    }
    try {
      return Q + "";
    } catch ($) {
    }
  }
  return "";
}
var F5 = Jq;
var Gq = /[\\^$.*+?()[\]{}|]/g;
var Wq = /^\[object .+?Constructor\]$/;
var Hq = Function.prototype;
var Bq = Object.prototype;
var zq = Hq.toString;
var Kq = Bq.hasOwnProperty;
var Vq = RegExp("^" + zq.call(Kq).replace(Gq, "\\$&").replace(/hasOwnProperty|(function).*?(?=\\\()| for .+?(?=\\\])/g, "$1.*?") + "$");
function qq(Q) {
  if (!p9(Q) || L5(Q)) return false;
  var $ = q5(Q) ? Vq : Wq;
  return $.test(F5(Q));
}
var N5 = qq;
function Uq(Q, $) {
  return Q == null ? void 0 : Q[$];
}
var O5 = Uq;
function Lq(Q, $) {
  var X = O5(Q, $);
  return N5(X) ? X : void 0;
}
var i9 = Lq;
var Fq = i9(Object, "create");
var b1 = Fq;
function Nq() {
  this.__data__ = b1 ? b1(null) : {}, this.size = 0;
}
var D5 = Nq;
function Oq(Q) {
  var $ = this.has(Q) && delete this.__data__[Q];
  return this.size -= $ ? 1 : 0, $;
}
var w5 = Oq;
var Dq = "__lodash_hash_undefined__";
var wq = Object.prototype;
var Mq = wq.hasOwnProperty;
function Aq(Q) {
  var $ = this.__data__;
  if (b1) {
    var X = $[Q];
    return X === Dq ? void 0 : X;
  }
  return Mq.call($, Q) ? $[Q] : void 0;
}
var M5 = Aq;
var jq = Object.prototype;
var Rq = jq.hasOwnProperty;
function Iq(Q) {
  var $ = this.__data__;
  return b1 ? $[Q] !== void 0 : Rq.call($, Q);
}
var A5 = Iq;
var Eq = "__lodash_hash_undefined__";
function Pq(Q, $) {
  var X = this.__data__;
  return this.size += this.has(Q) ? 0 : 1, X[Q] = b1 && $ === void 0 ? Eq : $, this;
}
var j5 = Pq;
function m6(Q) {
  var $ = -1, X = Q == null ? 0 : Q.length;
  this.clear();
  while (++$ < X) {
    var Y = Q[$];
    this.set(Y[0], Y[1]);
  }
}
m6.prototype.clear = D5;
m6.prototype.delete = w5;
m6.prototype.get = M5;
m6.prototype.has = A5;
m6.prototype.set = j5;
var mQ = m6;
function bq() {
  this.__data__ = [], this.size = 0;
}
var R5 = bq;
function Zq(Q, $) {
  return Q === $ || Q !== Q && $ !== $;
}
var I5 = Zq;
function Cq(Q, $) {
  var X = Q.length;
  while (X--) if (I5(Q[X][0], $)) return X;
  return -1;
}
var h1 = Cq;
var Sq = Array.prototype;
var _q = Sq.splice;
function kq(Q) {
  var $ = this.__data__, X = h1($, Q);
  if (X < 0) return false;
  var Y = $.length - 1;
  if (X == Y) $.pop();
  else _q.call($, X, 1);
  return --this.size, true;
}
var E5 = kq;
function vq(Q) {
  var $ = this.__data__, X = h1($, Q);
  return X < 0 ? void 0 : $[X][1];
}
var P5 = vq;
function Tq(Q) {
  return h1(this.__data__, Q) > -1;
}
var b5 = Tq;
function xq(Q, $) {
  var X = this.__data__, Y = h1(X, Q);
  if (Y < 0) ++this.size, X.push([Q, $]);
  else X[Y][1] = $;
  return this;
}
var Z5 = xq;
function l6(Q) {
  var $ = -1, X = Q == null ? 0 : Q.length;
  this.clear();
  while (++$ < X) {
    var Y = Q[$];
    this.set(Y[0], Y[1]);
  }
}
l6.prototype.clear = R5;
l6.prototype.delete = E5;
l6.prototype.get = P5;
l6.prototype.has = b5;
l6.prototype.set = Z5;
var C5 = l6;
var yq = i9(f6, "Map");
var S5 = yq;
function gq() {
  this.size = 0, this.__data__ = { hash: new mQ(), map: new (S5 || C5)(), string: new mQ() };
}
var _5 = gq;
function hq(Q) {
  var $ = typeof Q;
  return $ == "string" || $ == "number" || $ == "symbol" || $ == "boolean" ? Q !== "__proto__" : Q === null;
}
var k5 = hq;
function fq(Q, $) {
  var X = Q.__data__;
  return k5($) ? X[typeof $ == "string" ? "string" : "hash"] : X.map;
}
var f1 = fq;
function uq(Q) {
  var $ = f1(this, Q).delete(Q);
  return this.size -= $ ? 1 : 0, $;
}
var v5 = uq;
function mq(Q) {
  return f1(this, Q).get(Q);
}
var T5 = mq;
function lq(Q) {
  return f1(this, Q).has(Q);
}
var x5 = lq;
function cq(Q, $) {
  var X = f1(this, Q), Y = X.size;
  return X.set(Q, $), this.size += X.size == Y ? 0 : 1, this;
}
var y5 = cq;
function c6(Q) {
  var $ = -1, X = Q == null ? 0 : Q.length;
  this.clear();
  while (++$ < X) {
    var Y = Q[$];
    this.set(Y[0], Y[1]);
  }
}
c6.prototype.clear = _5;
c6.prototype.delete = v5;
c6.prototype.get = T5;
c6.prototype.has = x5;
c6.prototype.set = y5;
var lQ = c6;
var pq = "Expected a function";
function cQ(Q, $) {
  if (typeof Q != "function" || $ != null && typeof $ != "function") throw TypeError(pq);
  var X = function() {
    var Y = arguments, J = $ ? $.apply(this, Y) : Y[0], G = X.cache;
    if (G.has(J)) return G.get(J);
    var W = Q.apply(this, Y);
    return X.cache = G.set(J, W) || G, W;
  };
  return X.cache = new (cQ.Cache || lQ)(), X;
}
cQ.Cache = lQ;
var $1 = cQ;
var p6 = $1(() => {
  return (process.env.CLAUDE_CONFIG_DIR ?? dq(iq(), ".claude")).normalize("NFC");
}, () => process.env.CLAUDE_CONFIG_DIR);
function x4(Q) {
  if (!Q) return false;
  if (typeof Q === "boolean") return Q;
  let $ = Q.toLowerCase().trim();
  return ["1", "true", "yes", "on"].includes($);
}
var i6;
var d6 = null;
function tq() {
  if (d6) return d6;
  if (!process.env.DEBUG_CLAUDE_AGENT_SDK) return i6 = null, d6 = Promise.resolve(), d6;
  let Q = g5(p6(), "debug");
  return i6 = g5(Q, `sdk-${nq()}.txt`), process.stderr.write(`SDK debug logs: ${i6}
`), d6 = rq(Q, { recursive: true }).then(() => {
  }).catch(() => {
  }), d6;
}
function i0(Q) {
  if (i6 === null) return;
  let X = `${(/* @__PURE__ */ new Date()).toISOString()} ${Q}
`;
  tq().then(() => {
    if (i6) oq(i6, X).catch(() => {
    });
  });
}
function sq() {
  let Q = "";
  if (typeof process < "u" && typeof process.cwd === "function" && typeof h5 === "function") {
    let X = aq();
    try {
      Q = h5(X).normalize("NFC");
    } catch {
      Q = X.normalize("NFC");
    }
  }
  return { originalCwd: Q, projectRoot: Q, totalCostUSD: 0, totalAPIDuration: 0, totalAPIDurationWithoutRetries: 0, totalToolDuration: 0, tokenSaverBytesSaved: 0, tokenSaverHits: 0, turnHookDurationMs: 0, turnToolDurationMs: 0, turnClassifierDurationMs: 0, turnToolCount: 0, turnHookCount: 0, turnClassifierCount: 0, startTime: Date.now(), lastInteractionTime: Date.now(), totalLinesAdded: 0, totalLinesRemoved: 0, hasUnknownModelCost: false, cwd: Q, modelUsage: {}, mainLoopModelOverride: void 0, initialMainLoopModel: null, modelStrings: null, isInteractive: false, kairosActive: false, sdkAgentProgressSummariesEnabled: false, userMsgOptIn: false, clientType: "cli", sessionSource: void 0, questionPreviewFormat: void 0, sessionIngressToken: void 0, oauthTokenFromFd: void 0, apiKeyFromFd: void 0, flagSettingsPath: void 0, flagSettingsInline: null, allowedSettingSources: ["userSettings", "projectSettings", "localSettings", "flagSettings", "policySettings"], meter: null, sessionCounter: null, locCounter: null, prCounter: null, commitCounter: null, costCounter: null, tokenCounter: null, codeEditToolDecisionCounter: null, activeTimeCounter: null, statsStore: null, sessionId: n9(), parentSessionId: void 0, loggerProvider: null, eventLogger: null, meterProvider: null, tracerProvider: null, agentColorMap: /* @__PURE__ */ new Map(), agentColorIndex: 0, lastAPIRequest: null, lastAPIRequestMessages: null, lastClassifierRequests: null, cachedClaudeMdContent: null, inMemoryErrorLog: [], inlinePlugins: [], chromeFlagOverride: void 0, useCoworkPlugins: false, sessionBypassPermissionsMode: false, scheduledTasksEnabled: false, sessionCronTasks: [], sessionCreatedTeams: /* @__PURE__ */ new Set(), sessionTrustAccepted: false, sessionPersistenceDisabled: false, hasExitedPlanMode: false, needsPlanModeExitAttachment: false, needsAutoModeExitAttachment: false, lspRecommendationShownThisSession: false, initJsonSchema: null, registeredHooks: null, planSlugCache: /* @__PURE__ */ new Map(), teleportedSessionInfo: null, invokedSkills: /* @__PURE__ */ new Map(), slowOperations: [], sdkBetas: void 0, mainThreadAgentType: void 0, isRemoteMode: false, isInWorktree: false, ...{}, directConnectServerUrl: void 0, systemPromptSectionCache: /* @__PURE__ */ new Map(), lastEmittedDate: null, additionalDirectoriesForClaudeMd: [], allowedChannels: [], hasDevChannels: false, sessionProjectDir: null, promptCache1hAllowlist: null, promptId: null, lastMainRequestId: void 0, lastApiCompletionTimestamp: null, pendingPostCompaction: false };
}
var eq = sq();
function f5() {
  return eq.sessionId;
}
function u5({ writeFn: Q, flushIntervalMs: $ = 1e3, maxBufferSize: X = 100, maxBufferBytes: Y = 1 / 0, immediateMode: J = false }) {
  let G = [], W = 0, H = null, B = null;
  function z() {
    if (H) clearTimeout(H), H = null;
  }
  function K() {
    if (B) Q(B.join("")), B = null;
    if (G.length === 0) return;
    Q(G.join("")), G = [], W = 0, z();
  }
  function U() {
    if (!H) H = setTimeout(K, $);
  }
  function q() {
    if (B) {
      B.push(...G), G = [], W = 0, z();
      return;
    }
    let V = G;
    G = [], W = 0, z(), B = V, setImmediate(() => {
      let L = B;
      if (B = null, L) Q(L.join(""));
    });
  }
  return { write(V) {
    if (J) {
      Q(V);
      return;
    }
    if (G.push(V), W += V.length, U(), G.length >= X || W >= Y) q();
  }, flush: K, dispose() {
    K();
  } };
}
var m5 = /* @__PURE__ */ new Set();
function l5(Q) {
  return m5.add(Q), () => m5.delete(Q);
}
var c5 = $1((Q) => {
  if (!Q || Q.trim() === "") return null;
  let $ = Q.split(",").map((G) => G.trim()).filter(Boolean);
  if ($.length === 0) return null;
  let X = $.some((G) => G.startsWith("!")), Y = $.some((G) => !G.startsWith("!"));
  if (X && Y) return null;
  let J = $.map((G) => G.replace(/^!/, "").toLowerCase());
  return { include: X ? [] : J, exclude: X ? J : [], isExclusive: X };
});
function QU(Q) {
  let $ = [], X = Q.match(/^MCP server ["']([^"']+)["']/);
  if (X && X[1]) $.push("mcp"), $.push(X[1].toLowerCase());
  else {
    let G = Q.match(/^([^:[]+):/);
    if (G && G[1]) $.push(G[1].trim().toLowerCase());
  }
  let Y = Q.match(/^\[([^\]]+)]/);
  if (Y && Y[1]) $.push(Y[1].trim().toLowerCase());
  if (Q.toLowerCase().includes("1p event:")) $.push("1p");
  let J = Q.match(/:\s*([^:]+?)(?:\s+(?:type|mode|status|event))?:/);
  if (J && J[1]) {
    let G = J[1].trim().toLowerCase();
    if (G.length < 30 && !G.includes(" ")) $.push(G);
  }
  return Array.from(new Set($));
}
function $U(Q, $) {
  if (!$) return true;
  if (Q.length === 0) return false;
  if ($.isExclusive) return !Q.some((X) => $.exclude.includes(X));
  else return Q.some((X) => $.include.includes(X));
}
function p5(Q, $) {
  if (!$) return true;
  let X = QU(Q);
  return $U(X, $);
}
var KU = { cwd() {
  return process.cwd();
}, existsSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.existsSync(${Q})`, 0);
    return u.existsSync(Q);
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, async stat(Q) {
  return XU(Q);
}, async readdir(Q) {
  return YU(Q, { withFileTypes: true });
}, async unlink(Q) {
  return JU(Q);
}, async rmdir(Q) {
  return GU(Q);
}, async rm(Q, $) {
  return WU(Q, $);
}, async mkdir(Q, $) {
  try {
    await HU(Q, { recursive: true, ...$ });
  } catch (X) {
    if (X.code !== "EEXIST") throw X;
  }
}, async readFile(Q, $) {
  return d5(Q, { encoding: $.encoding });
}, async rename(Q, $) {
  return BU(Q, $);
}, statSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.statSync(${Q})`, 0);
    return u.statSync(Q);
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, lstatSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.lstatSync(${Q})`, 0);
    return u.lstatSync(Q);
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, readFileSync(Q, $) {
  let Y = [];
  try {
    const X = $0(Y, q0`fs.readFileSync(${Q})`, 0);
    return u.readFileSync(Q, { encoding: $.encoding });
  } catch (J) {
    var G = J, W = 1;
  } finally {
    X0(Y, G, W);
  }
}, readFileBytesSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.readFileBytesSync(${Q})`, 0);
    return u.readFileSync(Q);
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, readSync(Q, $) {
  let J = [];
  try {
    const X = $0(J, q0`fs.readSync(${Q}, ${$.length} bytes)`, 0);
    let Y = void 0;
    try {
      Y = u.openSync(Q, "r");
      let B = Buffer.alloc($.length), z = u.readSync(Y, B, 0, $.length, 0);
      return { buffer: B, bytesRead: z };
    } finally {
      if (Y) u.closeSync(Y);
    }
  } catch (G) {
    var W = G, H = 1;
  } finally {
    X0(J, W, H);
  }
}, appendFileSync(Q, $, X) {
  let J = [];
  try {
    const Y = $0(J, q0`fs.appendFileSync(${Q}, ${$.length} chars)`, 0);
    if (X?.mode !== void 0) try {
      let B = u.openSync(Q, "ax", X.mode);
      try {
        u.appendFileSync(B, $);
      } finally {
        u.closeSync(B);
      }
      return;
    } catch (B) {
      if (B.code !== "EEXIST") throw B;
    }
    u.appendFileSync(Q, $);
  } catch (G) {
    var W = G, H = 1;
  } finally {
    X0(J, W, H);
  }
}, copyFileSync(Q, $) {
  let Y = [];
  try {
    const X = $0(Y, q0`fs.copyFileSync(${Q} → ${$})`, 0);
    u.copyFileSync(Q, $);
  } catch (J) {
    var G = J, W = 1;
  } finally {
    X0(Y, G, W);
  }
}, unlinkSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.unlinkSync(${Q})`, 0);
    u.unlinkSync(Q);
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, renameSync(Q, $) {
  let Y = [];
  try {
    const X = $0(Y, q0`fs.renameSync(${Q} → ${$})`, 0);
    u.renameSync(Q, $);
  } catch (J) {
    var G = J, W = 1;
  } finally {
    X0(Y, G, W);
  }
}, linkSync(Q, $) {
  let Y = [];
  try {
    const X = $0(Y, q0`fs.linkSync(${Q} → ${$})`, 0);
    u.linkSync(Q, $);
  } catch (J) {
    var G = J, W = 1;
  } finally {
    X0(Y, G, W);
  }
}, symlinkSync(Q, $, X) {
  let J = [];
  try {
    const Y = $0(J, q0`fs.symlinkSync(${Q} → ${$})`, 0);
    u.symlinkSync(Q, $, X);
  } catch (G) {
    var W = G, H = 1;
  } finally {
    X0(J, W, H);
  }
}, readlinkSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.readlinkSync(${Q})`, 0);
    return u.readlinkSync(Q);
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, realpathSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.realpathSync(${Q})`, 0);
    return u.realpathSync(Q).normalize("NFC");
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, mkdirSync(Q, $) {
  let J = [];
  try {
    const X = $0(J, q0`fs.mkdirSync(${Q})`, 0);
    let Y = { recursive: true };
    if ($?.mode !== void 0) Y.mode = $.mode;
    try {
      u.mkdirSync(Q, Y);
    } catch (B) {
      if (B.code !== "EEXIST") throw B;
    }
  } catch (G) {
    var W = G, H = 1;
  } finally {
    X0(J, W, H);
  }
}, readdirSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.readdirSync(${Q})`, 0);
    return u.readdirSync(Q, { withFileTypes: true });
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, readdirStringSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.readdirStringSync(${Q})`, 0);
    return u.readdirSync(Q);
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, isDirEmptySync(Q) {
  let Y = [];
  try {
    const $ = $0(Y, q0`fs.isDirEmptySync(${Q})`, 0);
    let X = this.readdirSync(Q);
    return X.length === 0;
  } catch (J) {
    var G = J, W = 1;
  } finally {
    X0(Y, G, W);
  }
}, rmdirSync(Q) {
  let X = [];
  try {
    const $ = $0(X, q0`fs.rmdirSync(${Q})`, 0);
    u.rmdirSync(Q);
  } catch (Y) {
    var J = Y, G = 1;
  } finally {
    X0(X, J, G);
  }
}, rmSync(Q, $) {
  let Y = [];
  try {
    const X = $0(Y, q0`fs.rmSync(${Q})`, 0);
    u.rmSync(Q, $);
  } catch (J) {
    var G = J, W = 1;
  } finally {
    X0(Y, G, W);
  }
}, createWriteStream(Q) {
  return u.createWriteStream(Q);
}, async readFileBytes(Q, $) {
  if ($ === void 0) return d5(Q);
  let X = await zU(Q, "r");
  try {
    let { size: Y } = await X.stat(), J = Math.min(Y, $), G = Buffer.allocUnsafe(J), W = 0;
    while (W < J) {
      let { bytesRead: H } = await X.read(G, W, J - W, W);
      if (H === 0) break;
      W += H;
    }
    return W < J ? G.subarray(0, W) : G;
  } finally {
    await X.close();
  }
} };
var VU = KU;
function pQ() {
  return VU;
}
function qU(Q, $) {
  if (Q.destroyed) return;
  Q.write($);
}
function i5(Q) {
  qU(process.stderr, Q);
}
var iQ = { verbose: 0, debug: 1, info: 2, warn: 3, error: 4 };
var OU = $1(() => {
  let Q = process.env.CLAUDE_CODE_DEBUG_LOG_LEVEL?.toLowerCase().trim();
  if (Q && Object.hasOwn(iQ, Q)) return Q;
  return "debug";
});
var DU = false;
var nQ = $1(() => {
  return DU || x4(process.env.DEBUG) || x4(process.env.DEBUG_SDK) || process.argv.includes("--debug") || process.argv.includes("-d") || r5() || process.argv.some((Q) => Q.startsWith("--debug=")) || t5() !== null;
});
var wU = $1(() => {
  let Q = process.argv.find((X) => X.startsWith("--debug="));
  if (!Q) return null;
  let $ = Q.substring(8);
  return c5($);
});
var r5 = $1(() => {
  return process.argv.includes("--debug-to-stderr") || process.argv.includes("-d2e");
});
var t5 = $1(() => {
  for (let Q = 0; Q < process.argv.length; Q++) {
    let $ = process.argv[Q];
    if ($.startsWith("--debug-file=")) return $.substring(13);
    if ($ === "--debug-file" && Q + 1 < process.argv.length) return process.argv[Q + 1];
  }
  return null;
});
function MU(Q) {
  if (!nQ()) return false;
  if (typeof process > "u" || typeof process.versions > "u" || typeof process.versions.node > "u") return false;
  let $ = wU();
  return p5(Q, $);
}
var AU = false;
var o9 = null;
var dQ = Promise.resolve();
async function jU(Q, $, X, Y) {
  if (Q) await LU($, { recursive: true }).catch(() => {
  });
  await UU(X, Y), s5();
}
function RU() {
}
function IU() {
  if (!o9) {
    let Q = null;
    o9 = u5({ writeFn: ($) => {
      let X = a5(), Y = n5(X), J = Q !== Y;
      if (Q = Y, nQ()) {
        if (J) try {
          pQ().mkdirSync(Y);
        } catch {
        }
        pQ().appendFileSync(X, $), s5();
        return;
      }
      dQ = dQ.then(jU.bind(null, J, Y, X, $)).catch(RU);
    }, flushIntervalMs: 1e3, maxBufferSize: 100, immediateMode: nQ() }), l5(async () => {
      o9?.dispose(), await dQ;
    });
  }
  return o9;
}
function N1(Q, { level: $ } = { level: "debug" }) {
  if (iQ[$] < iQ[OU()]) return;
  if (!MU(Q)) return;
  if (AU && Q.includes(`
`)) Q = W0(Q);
  let Y = `${(/* @__PURE__ */ new Date()).toISOString()} [${$.toUpperCase()}] ${Q.trim()}
`;
  if (r5()) {
    i5(Y);
    return;
  }
  IU().write(Y);
}
function a5() {
  return t5() ?? process.env.CLAUDE_CODE_DEBUG_LOGS_DIR ?? o5(p6(), "debug", `${f5()}.txt`);
}
var s5 = $1(async () => {
  try {
    let Q = a5(), $ = n5(Q), X = o5($, "latest");
    await FU(X).catch(() => {
    }), await NU(Q, X);
  } catch {
  }
});
var PC = (() => {
  let Q = process.env.CLAUDE_CODE_SLOW_OPERATION_THRESHOLD_MS;
  if (Q !== void 0) {
    let $ = Number(Q);
    if (!Number.isNaN($) && $ >= 0) return $;
  }
  return 1 / 0;
})();
var EU = { [Symbol.dispose]() {
} };
function PU() {
  return EU;
}
var q0 = PU;
function W0(Q, $, X) {
  let J = [];
  try {
    const Y = $0(J, q0`JSON.stringify(${Q})`, 0);
    return JSON.stringify(Q, $, X);
  } catch (G) {
    var W = G, H = 1;
  } finally {
    X0(J, W, H);
  }
}
var O1 = (Q, $) => {
  let Y = [];
  try {
    const X = $0(Y, q0`JSON.parse(${Q})`, 0);
    return typeof $ > "u" ? JSON.parse(Q) : JSON.parse(Q, $);
  } catch (J) {
    var G = J, W = 1;
  } finally {
    X0(Y, G, W);
  }
};
function bU(Q) {
  let $ = Q.trim();
  return $.startsWith("{") && $.endsWith("}");
}
function e5(Q, $) {
  let X = { ...Q };
  if ($) {
    let Y = X.settings;
    if (Y && !bU(Y)) throw Error("Cannot use both a settings file path and the sandbox option. Include the sandbox configuration in your settings file instead.");
    let J = { sandbox: $ };
    if (Y) try {
      J = { ...O1(Y), sandbox: $ };
    } catch {
    }
    X.settings = W0(J);
  }
  return X;
}
var SU = 2e3;
var y4 = class {
  options;
  process;
  processStdin;
  processStdout;
  ready = false;
  abortController;
  exitError;
  exitListeners = [];
  processExitHandler;
  abortHandler;
  constructor(Q) {
    this.options = Q;
    this.abortController = Q.abortController || g6(), this.initialize();
  }
  getDefaultExecutable() {
    return h6() ? "bun" : "node";
  }
  spawnLocalProcess(Q) {
    let { command: $, args: X, cwd: Y, env: J, signal: G } = Q, W = J.DEBUG_CLAUDE_AGENT_SDK || this.options.stderr ? "pipe" : "ignore", H = ZU($, X, { cwd: Y, stdio: ["pipe", "pipe", W], signal: G, env: J, windowsHide: true });
    if (J.DEBUG_CLAUDE_AGENT_SDK || this.options.stderr) H.stderr.on("data", (z) => {
      let K = z.toString();
      if (i0(K), this.options.stderr) this.options.stderr(K);
    });
    return { stdin: H.stdin, stdout: H.stdout, get killed() {
      return H.killed;
    }, get exitCode() {
      return H.exitCode;
    }, kill: H.kill.bind(H), on: H.on.bind(H), once: H.once.bind(H), off: H.off.bind(H) };
  }
  initialize() {
    try {
      let { additionalDirectories: Q = [], agent: $, betas: X, cwd: Y, executable: J = this.getDefaultExecutable(), executableArgs: G = [], extraArgs: W = {}, pathToClaudeCodeExecutable: H, env: B = { ...process.env }, thinkingConfig: z, maxTurns: K, maxBudgetUsd: U, model: q, fallbackModel: V, jsonSchema: L, permissionMode: F, allowDangerouslySkipPermissions: w, permissionPromptToolName: D, continueConversation: M, resume: R, settingSources: Z, allowedTools: v = [], disallowedTools: O0 = [], tools: D0, mcpServers: d0, strictMcpConfig: B6, canUseTool: F1, includePartialMessages: z6, plugins: y1, sandbox: K6 } = this.options, h = ["--output-format", "stream-json", "--verbose", "--input-format", "stream-json"];
      if (z) switch (z.type) {
        case "enabled":
          if (z.budgetTokens === void 0) h.push("--thinking", "adaptive");
          else h.push("--max-thinking-tokens", z.budgetTokens.toString());
          break;
        case "disabled":
          h.push("--thinking", "disabled");
          break;
        case "adaptive":
          h.push("--thinking", "adaptive");
          break;
      }
      if (this.options.effort) h.push("--effort", this.options.effort);
      if (K) h.push("--max-turns", K.toString());
      if (U !== void 0) h.push("--max-budget-usd", U.toString());
      if (q) h.push("--model", q);
      if ($) h.push("--agent", $);
      if (X && X.length > 0) h.push("--betas", X.join(","));
      if (L) h.push("--json-schema", W0(L));
      if (this.options.debugFile) h.push("--debug-file", this.options.debugFile);
      else if (this.options.debug) h.push("--debug");
      if (B.DEBUG_CLAUDE_AGENT_SDK) h.push("--debug-to-stderr");
      if (F1) {
        if (D) throw Error("canUseTool callback cannot be used with permissionPromptToolName. Please use one or the other.");
        h.push("--permission-prompt-tool", "stdio");
      } else if (D) h.push("--permission-prompt-tool", D);
      if (M) h.push("--continue");
      if (R) h.push("--resume", R);
      if (this.options.proactive) h.push("--proactive");
      if (this.options.assistant) h.push("--assistant");
      if (v.length > 0) h.push("--allowedTools", v.join(","));
      if (O0.length > 0) h.push("--disallowedTools", O0.join(","));
      if (D0 !== void 0) if (Array.isArray(D0)) if (D0.length === 0) h.push("--tools", "");
      else h.push("--tools", D0.join(","));
      else h.push("--tools", "default");
      if (d0 && Object.keys(d0).length > 0) h.push("--mcp-config", W0({ mcpServers: d0 }));
      if (Z) h.push("--setting-sources", Z.join(","));
      if (B6) h.push("--strict-mcp-config");
      if (F) h.push("--permission-mode", F);
      if (w) h.push("--allow-dangerously-skip-permissions");
      if (V) {
        if (q && V === q) throw Error("Fallback model cannot be the same as the main model. Please specify a different model for fallbackModel option.");
        h.push("--fallback-model", V);
      }
      if (z6) h.push("--include-partial-messages");
      for (let C0 of Q) h.push("--add-dir", C0);
      if (y1 && y1.length > 0) for (let C0 of y1) if (C0.type === "local") h.push("--plugin-dir", C0.path);
      else throw Error(`Unsupported plugin type: ${C0.type}`);
      if (this.options.forkSession) h.push("--fork-session");
      if (this.options.resumeSessionAt) h.push("--resume-session-at", this.options.resumeSessionAt);
      if (this.options.sessionId) h.push("--session-id", this.options.sessionId);
      if (this.options.persistSession === false) h.push("--no-session-persistence");
      let S4 = { ...W ?? {} };
      if (this.options.settings) S4.settings = this.options.settings;
      let gQ = e5(S4, K6);
      for (let [C0, g1] of Object.entries(gQ)) if (g1 === null) h.push(`--${C0}`);
      else h.push(`--${C0}`, g1);
      if (!B.CLAUDE_CODE_ENTRYPOINT) B.CLAUDE_CODE_ENTRYPOINT = "sdk-ts";
      if (delete B.NODE_OPTIONS, B.DEBUG_CLAUDE_AGENT_SDK) B.DEBUG = "1";
      else delete B.DEBUG;
      let _4 = _U(H), k4 = _4 ? H : J, V6 = _4 ? [...G, ...h] : [...G, H, ...h], c9 = { command: k4, args: V6, cwd: Y, env: B, signal: this.abortController.signal };
      if (this.options.spawnClaudeCodeProcess) i0(`Spawning Claude Code (custom): ${k4} ${V6.join(" ")}`), this.process = this.options.spawnClaudeCodeProcess(c9);
      else i0(`Spawning Claude Code: ${k4} ${V6.join(" ")}`), this.process = this.spawnLocalProcess(c9);
      this.processStdin = this.process.stdin, this.processStdout = this.process.stdout;
      let v6 = () => {
        if (this.process && !this.process.killed) this.process.kill("SIGTERM");
      };
      this.processExitHandler = v6, this.abortHandler = v6, process.on("exit", this.processExitHandler), this.abortController.signal.addEventListener("abort", this.abortHandler), this.process.on("error", (C0) => {
        if (this.ready = false, this.abortController.signal.aborted) this.exitError = new T0("Claude Code process aborted by user");
        else if (C0.code === "ENOENT") {
          let g1 = _4 ? `Claude Code native binary not found at ${H}. Please ensure Claude Code is installed via native installer or specify a valid path with options.pathToClaudeCodeExecutable.` : `Claude Code executable not found at ${H}. Is options.pathToClaudeCodeExecutable set?`;
          this.exitError = ReferenceError(g1), i0(this.exitError.message);
        } else this.exitError = Error(`Failed to spawn Claude Code process: ${C0.message}`), i0(this.exitError.message);
      }), this.process.on("exit", (C0, g1) => {
        if (this.ready = false, this.abortController.signal.aborted) this.exitError = new T0("Claude Code process aborted by user");
        else {
          let T6 = this.getProcessExitError(C0, g1);
          if (T6) this.exitError = T6, i0(T6.message);
        }
      }), this.ready = true;
    } catch (Q) {
      throw this.ready = false, Q;
    }
  }
  getProcessExitError(Q, $) {
    if (Q !== 0 && Q !== null) return Error(`Claude Code process exited with code ${Q}`);
    else if ($) return Error(`Claude Code process terminated by signal ${$}`);
    return;
  }
  write(Q) {
    if (this.abortController.signal.aborted) throw new T0("Operation aborted");
    if (!this.ready || !this.processStdin) throw Error("ProcessTransport is not ready for writing");
    if (this.process?.killed || this.process?.exitCode !== null) throw Error("Cannot write to terminated process");
    if (this.exitError) throw Error(`Cannot write to process that exited with error: ${this.exitError.message}`);
    i0(`[ProcessTransport] Writing to stdin: ${Q.substring(0, 100)}`);
    try {
      if (!this.processStdin.write(Q)) i0("[ProcessTransport] Write buffer full, data queued");
    } catch ($) {
      throw this.ready = false, Error(`Failed to write to process stdin: ${$.message}`);
    }
  }
  close() {
    if (this.processStdin) this.processStdin.end(), this.processStdin = void 0;
    if (this.abortHandler) this.abortController.signal.removeEventListener("abort", this.abortHandler), this.abortHandler = void 0;
    for (let { handler: $ } of this.exitListeners) this.process?.off("exit", $);
    this.exitListeners = [];
    let Q = this.process;
    if (Q && !Q.killed && Q.exitCode === null) setTimeout(($) => {
      if ($.killed || $.exitCode !== null) return;
      $.kill("SIGTERM"), setTimeout((X) => {
        if (X.exitCode === null) X.kill("SIGKILL");
      }, 5e3, $).unref();
    }, SU, Q).unref(), Q.once("exit", () => {
      if (this.processExitHandler) process.off("exit", this.processExitHandler), this.processExitHandler = void 0;
    });
    else if (this.processExitHandler) process.off("exit", this.processExitHandler), this.processExitHandler = void 0;
    this.ready = false;
  }
  isReady() {
    return this.ready;
  }
  async *readMessages() {
    if (!this.processStdout) throw Error("ProcessTransport output stream not available");
    let Q = CU({ input: this.processStdout });
    try {
      for await (let $ of Q) if ($.trim()) try {
        yield O1($);
      } catch (X) {
        throw i0(`Non-JSON stdout: ${$}`), Error(`CLI output was not valid JSON. This may indicate an error during startup. Output: ${$.slice(0, 200)}${$.length > 200 ? "..." : ""}`);
      }
      await this.waitForExit();
    } catch ($) {
      throw $;
    } finally {
      Q.close();
    }
  }
  endInput() {
    if (this.processStdin) this.processStdin.end();
  }
  getInputStream() {
    return this.processStdin;
  }
  onExit(Q) {
    if (!this.process) return () => {
    };
    let $ = (X, Y) => {
      let J = this.getProcessExitError(X, Y);
      Q(J);
    };
    return this.process.on("exit", $), this.exitListeners.push({ callback: Q, handler: $ }), () => {
      if (this.process) this.process.off("exit", $);
      let X = this.exitListeners.findIndex((Y) => Y.handler === $);
      if (X !== -1) this.exitListeners.splice(X, 1);
    };
  }
  async waitForExit() {
    if (!this.process) {
      if (this.exitError) throw this.exitError;
      return;
    }
    if (this.process.exitCode !== null || this.process.killed) {
      if (this.exitError) throw this.exitError;
      return;
    }
    return new Promise((Q, $) => {
      let X = (J, G) => {
        if (this.abortController.signal.aborted) {
          $(new T0("Operation aborted"));
          return;
        }
        let W = this.getProcessExitError(J, G);
        if (W) $(W);
        else Q();
      };
      this.process.once("exit", X);
      let Y = (J) => {
        this.process.off("exit", X), $(J);
      };
      this.process.once("error", Y), this.process.once("exit", () => {
        this.process.off("error", Y);
      });
    });
  }
};
function _U(Q) {
  return ![".js", ".mjs", ".tsx", ".ts", ".jsx"].some((X) => Q.endsWith(X));
}
var g4 = class {
  returned;
  queue = [];
  readResolve;
  readReject;
  isDone = false;
  hasError;
  started = false;
  constructor(Q) {
    this.returned = Q;
  }
  [Symbol.asyncIterator]() {
    if (this.started) throw Error("Stream can only be iterated once");
    return this.started = true, this;
  }
  next() {
    if (this.queue.length > 0) return Promise.resolve({ done: false, value: this.queue.shift() });
    if (this.isDone) return Promise.resolve({ done: true, value: void 0 });
    if (this.hasError) return Promise.reject(this.hasError);
    return new Promise((Q, $) => {
      this.readResolve = Q, this.readReject = $;
    });
  }
  enqueue(Q) {
    if (this.readResolve) {
      let $ = this.readResolve;
      this.readResolve = void 0, this.readReject = void 0, $({ done: false, value: Q });
    } else this.queue.push(Q);
  }
  done() {
    if (this.isDone = true, this.readResolve) {
      let Q = this.readResolve;
      this.readResolve = void 0, this.readReject = void 0, Q({ done: true, value: void 0 });
    }
  }
  error(Q) {
    if (this.hasError = Q, this.readReject) {
      let $ = this.readReject;
      this.readResolve = void 0, this.readReject = void 0, $(Q);
    }
  }
  return() {
    if (this.isDone = true, this.returned) this.returned();
    return Promise.resolve({ done: true, value: void 0 });
  }
};
var oQ = class {
  sendMcpMessage;
  isClosed = false;
  constructor(Q) {
    this.sendMcpMessage = Q;
  }
  onclose;
  onerror;
  onmessage;
  async start() {
  }
  async send(Q) {
    if (this.isClosed) throw Error("Transport is closed");
    this.sendMcpMessage(Q);
  }
  async close() {
    if (this.isClosed) return;
    this.isClosed = true, this.onclose?.();
  }
};
var h4 = class {
  transport;
  isSingleUserTurn;
  canUseTool;
  hooks;
  abortController;
  jsonSchema;
  initConfig;
  onElicitation;
  pendingControlResponses = /* @__PURE__ */ new Map();
  cleanupPerformed = false;
  sdkMessages;
  inputStream = new g4();
  initialization;
  cancelControllers = /* @__PURE__ */ new Map();
  hookCallbacks = /* @__PURE__ */ new Map();
  nextCallbackId = 0;
  sdkMcpTransports = /* @__PURE__ */ new Map();
  sdkMcpServerInstances = /* @__PURE__ */ new Map();
  pendingMcpResponses = /* @__PURE__ */ new Map();
  firstResultReceivedResolve;
  firstResultReceived = false;
  lastErrorResultText;
  hasBidirectionalNeeds() {
    return this.sdkMcpTransports.size > 0 || this.hooks !== void 0 && Object.keys(this.hooks).length > 0 || this.canUseTool !== void 0 || this.onElicitation !== void 0;
  }
  constructor(Q, $, X, Y, J, G = /* @__PURE__ */ new Map(), W, H, B) {
    this.transport = Q;
    this.isSingleUserTurn = $;
    this.canUseTool = X;
    this.hooks = Y;
    this.abortController = J;
    this.jsonSchema = W;
    this.initConfig = H;
    this.onElicitation = B;
    for (let [z, K] of G) this.connectSdkMcpServer(z, K);
    this.sdkMessages = this.readSdkMessages(), this.readMessages(), this.initialization = this.initialize(), this.initialization.catch(() => {
    });
  }
  setError(Q) {
    this.inputStream.error(Q);
  }
  async stopTask(Q) {
    await this.request({ subtype: "stop_task", task_id: Q });
  }
  close() {
    this.cleanup();
  }
  cleanup(Q) {
    if (this.cleanupPerformed) return;
    this.cleanupPerformed = true;
    try {
      this.transport.close();
      let $ = Error("Query closed before response received");
      for (let { reject: X } of this.pendingControlResponses.values()) X($);
      this.pendingControlResponses.clear();
      for (let { reject: X } of this.pendingMcpResponses.values()) X($);
      this.pendingMcpResponses.clear(), this.cancelControllers.clear(), this.hookCallbacks.clear();
      for (let X of this.sdkMcpTransports.values()) try {
        X.close();
      } catch {
      }
      if (this.sdkMcpTransports.clear(), Q) this.inputStream.error(Q);
      else this.inputStream.done();
    } catch ($) {
    }
  }
  next(...[Q]) {
    return this.sdkMessages.next(...[Q]);
  }
  return(Q) {
    return this.sdkMessages.return(Q);
  }
  throw(Q) {
    return this.sdkMessages.throw(Q);
  }
  [Symbol.asyncIterator]() {
    return this.sdkMessages;
  }
  [Symbol.asyncDispose]() {
    return this.sdkMessages[Symbol.asyncDispose]();
  }
  async readMessages() {
    try {
      for await (let Q of this.transport.readMessages()) {
        if (Q.type === "control_response") {
          let $ = this.pendingControlResponses.get(Q.response.request_id);
          if ($) $.handler(Q.response);
          continue;
        } else if (Q.type === "control_request") {
          this.handleControlRequest(Q);
          continue;
        } else if (Q.type === "control_cancel_request") {
          this.handleControlCancelRequest(Q);
          continue;
        } else if (Q.type === "keep_alive") continue;
        if (Q.type === "streamlined_text" || Q.type === "streamlined_tool_use_summary") continue;
        if (Q.type === "result") {
          if (this.lastErrorResultText = Q.is_error ? Q.subtype === "success" ? Q.result : Q.errors.join("; ") : void 0, this.firstResultReceived = true, this.firstResultReceivedResolve) this.firstResultReceivedResolve();
          if (this.isSingleUserTurn) N1("[Query.readMessages] First result received for single-turn query, closing stdin"), this.transport.endInput();
        } else this.lastErrorResultText = void 0;
        this.inputStream.enqueue(Q);
      }
      if (this.firstResultReceivedResolve) this.firstResultReceivedResolve();
      this.inputStream.done(), this.cleanup();
    } catch (Q) {
      if (this.firstResultReceivedResolve) this.firstResultReceivedResolve();
      if (this.lastErrorResultText !== void 0 && !(Q instanceof T0)) {
        let $ = Error(`Claude Code returned an error result: ${this.lastErrorResultText}`);
        N1(`[Query.readMessages] Replacing exit error with result text. Original: ${Q.message}`), this.inputStream.error($), this.cleanup($);
        return;
      }
      this.inputStream.error(Q), this.cleanup(Q);
    }
  }
  async handleControlRequest(Q) {
    let $ = new AbortController();
    this.cancelControllers.set(Q.request_id, $);
    try {
      let X = await this.processControlRequest(Q, $.signal), Y = { type: "control_response", response: { subtype: "success", request_id: Q.request_id, response: X } };
      await Promise.resolve(this.transport.write(W0(Y) + `
`));
    } catch (X) {
      let Y = { type: "control_response", response: { subtype: "error", request_id: Q.request_id, error: X.message || String(X) } };
      await Promise.resolve(this.transport.write(W0(Y) + `
`));
    } finally {
      this.cancelControllers.delete(Q.request_id);
    }
  }
  handleControlCancelRequest(Q) {
    let $ = this.cancelControllers.get(Q.request_id);
    if ($) $.abort(), this.cancelControllers.delete(Q.request_id);
  }
  async processControlRequest(Q, $) {
    if (Q.request.subtype === "can_use_tool") {
      if (!this.canUseTool) throw Error("canUseTool callback is not provided.");
      return { ...await this.canUseTool(Q.request.tool_name, Q.request.input, { signal: $, suggestions: Q.request.permission_suggestions, blockedPath: Q.request.blocked_path, decisionReason: Q.request.decision_reason, title: Q.request.title, displayName: Q.request.display_name, description: Q.request.description, toolUseID: Q.request.tool_use_id, agentID: Q.request.agent_id }), toolUseID: Q.request.tool_use_id };
    } else if (Q.request.subtype === "hook_callback") return await this.handleHookCallbacks(Q.request.callback_id, Q.request.input, Q.request.tool_use_id, $);
    else if (Q.request.subtype === "mcp_message") {
      let X = Q.request, Y = this.sdkMcpTransports.get(X.server_name);
      if (!Y) throw Error(`SDK MCP server not found: ${X.server_name}`);
      if ("method" in X.message && "id" in X.message && X.message.id !== null) return { mcp_response: await this.handleMcpControlRequest(X.server_name, X, Y) };
      else {
        if (Y.onmessage) Y.onmessage(X.message);
        return { mcp_response: { jsonrpc: "2.0", result: {}, id: 0 } };
      }
    } else if (Q.request.subtype === "elicitation") {
      let X = Q.request;
      if (this.onElicitation) return await this.onElicitation({ serverName: X.mcp_server_name, message: X.message, mode: X.mode, url: X.url, elicitationId: X.elicitation_id, requestedSchema: X.requested_schema }, { signal: $ });
      return { action: "decline" };
    }
    throw Error("Unsupported control request subtype: " + Q.request.subtype);
  }
  async *readSdkMessages() {
    for await (let Q of this.inputStream) yield Q;
  }
  async initialize() {
    let Q;
    if (this.hooks) {
      Q = {};
      for (let [J, G] of Object.entries(this.hooks)) if (G.length > 0) Q[J] = G.map((W) => {
        let H = [];
        for (let B of W.hooks) {
          let z = `hook_${this.nextCallbackId++}`;
          this.hookCallbacks.set(z, B), H.push(z);
        }
        return { matcher: W.matcher, hookCallbackIds: H, timeout: W.timeout };
      });
    }
    let $ = this.sdkMcpTransports.size > 0 ? Array.from(this.sdkMcpTransports.keys()) : void 0, X = { subtype: "initialize", hooks: Q, sdkMcpServers: $, jsonSchema: this.jsonSchema, systemPrompt: this.initConfig?.systemPrompt, appendSystemPrompt: this.initConfig?.appendSystemPrompt, agents: this.initConfig?.agents, promptSuggestions: this.initConfig?.promptSuggestions, agentProgressSummaries: this.initConfig?.agentProgressSummaries };
    return (await this.request(X)).response;
  }
  async interrupt() {
    await this.request({ subtype: "interrupt" });
  }
  async setPermissionMode(Q) {
    await this.request({ subtype: "set_permission_mode", mode: Q });
  }
  async setModel(Q) {
    await this.request({ subtype: "set_model", model: Q });
  }
  async setMaxThinkingTokens(Q) {
    await this.request({ subtype: "set_max_thinking_tokens", max_thinking_tokens: Q });
  }
  async applyFlagSettings(Q) {
    await this.request({ subtype: "apply_flag_settings", settings: Q });
  }
  async getSettings() {
    return (await this.request({ subtype: "get_settings" })).response;
  }
  async rewindFiles(Q, $) {
    return (await this.request({ subtype: "rewind_files", user_message_id: Q, dry_run: $?.dryRun })).response;
  }
  async cancelAsyncMessage(Q) {
    return (await this.request({ subtype: "cancel_async_message", message_uuid: Q })).response.cancelled;
  }
  async enableRemoteControl(Q) {
    return (await this.request({ subtype: "remote_control", enabled: Q })).response;
  }
  async setProactive(Q) {
    await this.request({ subtype: "set_proactive", enabled: Q });
  }
  async generateSessionTitle(Q, $) {
    return (await this.request({ subtype: "generate_session_title", description: Q, persist: $?.persist })).response.title;
  }
  async processPendingPermissionRequests(Q) {
    for (let $ of Q) if ($.request.subtype === "can_use_tool") this.handleControlRequest($).catch(() => {
    });
  }
  request(Q) {
    let $ = Math.random().toString(36).substring(2, 15), X = { request_id: $, type: "control_request", request: Q };
    return new Promise((Y, J) => {
      this.pendingControlResponses.set($, { handler: (G) => {
        if (this.pendingControlResponses.delete($), G.subtype === "success") Y(G);
        else if (J(Error(G.error)), G.pending_permission_requests) this.processPendingPermissionRequests(G.pending_permission_requests);
      }, reject: J }), Promise.resolve(this.transport.write(W0(X) + `
`));
    });
  }
  async initializationResult() {
    return this.initialization;
  }
  async supportedCommands() {
    return (await this.initialization).commands;
  }
  async supportedModels() {
    return (await this.initialization).models;
  }
  async supportedAgents() {
    return (await this.initialization).agents;
  }
  async reconnectMcpServer(Q) {
    await this.request({ subtype: "mcp_reconnect", serverName: Q });
  }
  async toggleMcpServer(Q, $) {
    await this.request({ subtype: "mcp_toggle", serverName: Q, enabled: $ });
  }
  async mcpAuthenticate(Q) {
    return (await this.request({ subtype: "mcp_authenticate", serverName: Q })).response;
  }
  async mcpClearAuth(Q) {
    return (await this.request({ subtype: "mcp_clear_auth", serverName: Q })).response;
  }
  async mcpSubmitOAuthCallbackUrl(Q, $) {
    return (await this.request({ subtype: "mcp_oauth_callback_url", serverName: Q, callbackUrl: $ })).response;
  }
  async claudeAuthenticate(Q) {
    return (await this.request({ subtype: "claude_authenticate", loginWithClaudeAi: Q })).response;
  }
  async claudeOAuthCallback(Q, $) {
    return (await this.request({ subtype: "claude_oauth_callback", authorizationCode: Q, state: $ })).response;
  }
  async claudeOAuthWaitForCompletion() {
    return (await this.request({ subtype: "claude_oauth_wait_for_completion" })).response;
  }
  async mcpServerStatus() {
    return (await this.request({ subtype: "mcp_status" })).response.mcpServers;
  }
  async setMcpServers(Q) {
    let $ = {}, X = {};
    for (let [H, B] of Object.entries(Q)) if (B.type === "sdk" && "instance" in B) $[H] = B.instance;
    else X[H] = B;
    let Y = new Set(this.sdkMcpServerInstances.keys()), J = new Set(Object.keys($));
    for (let H of Y) if (!J.has(H)) await this.disconnectSdkMcpServer(H);
    for (let [H, B] of Object.entries($)) if (!Y.has(H)) this.connectSdkMcpServer(H, B);
    let G = {};
    for (let H of Object.keys($)) G[H] = { type: "sdk", name: H };
    return (await this.request({ subtype: "mcp_set_servers", servers: { ...X, ...G } })).response;
  }
  async accountInfo() {
    return (await this.initialization).account;
  }
  async streamInput(Q) {
    N1("[Query.streamInput] Starting to process input stream");
    try {
      let $ = 0;
      for await (let X of Q) {
        if ($++, N1(`[Query.streamInput] Processing message ${$}: ${X.type}`), this.abortController?.signal.aborted) break;
        await Promise.resolve(this.transport.write(W0(X) + `
`));
      }
      if (N1(`[Query.streamInput] Finished processing ${$} messages from input stream`), $ > 0 && this.hasBidirectionalNeeds()) N1("[Query.streamInput] Has bidirectional needs, waiting for first result"), await this.waitForFirstResult();
      N1("[Query] Calling transport.endInput() to close stdin to CLI process"), this.transport.endInput();
    } catch ($) {
      if (!($ instanceof T0)) throw $;
    }
  }
  waitForFirstResult() {
    if (this.firstResultReceived) return N1("[Query.waitForFirstResult] Result already received, returning immediately"), Promise.resolve();
    return new Promise((Q) => {
      if (this.abortController?.signal.aborted) {
        Q();
        return;
      }
      this.abortController?.signal.addEventListener("abort", () => Q(), { once: true }), this.firstResultReceivedResolve = Q;
    });
  }
  handleHookCallbacks(Q, $, X, Y) {
    let J = this.hookCallbacks.get(Q);
    if (!J) throw Error(`No hook callback found for ID: ${Q}`);
    return J($, X, { signal: Y });
  }
  connectSdkMcpServer(Q, $) {
    let X = new oQ((Y) => this.sendMcpServerMessageToCli(Q, Y));
    this.sdkMcpTransports.set(Q, X), this.sdkMcpServerInstances.set(Q, $), $.connect(X);
  }
  async disconnectSdkMcpServer(Q) {
    let $ = this.sdkMcpTransports.get(Q);
    if ($) await $.close(), this.sdkMcpTransports.delete(Q);
    this.sdkMcpServerInstances.delete(Q);
  }
  sendMcpServerMessageToCli(Q, $) {
    if ("id" in $ && $.id !== null && $.id !== void 0) {
      let Y = `${Q}:${$.id}`, J = this.pendingMcpResponses.get(Y);
      if (J) {
        J.resolve($), this.pendingMcpResponses.delete(Y);
        return;
      }
    }
    let X = { type: "control_request", request_id: n9(), request: { subtype: "mcp_message", server_name: Q, message: $ } };
    this.transport.write(W0(X) + `
`);
  }
  handleMcpControlRequest(Q, $, X) {
    let Y = "id" in $.message ? $.message.id : null, J = `${Q}:${Y}`;
    return new Promise((G, W) => {
      let H = () => {
        this.pendingMcpResponses.delete(J);
      }, B = (K) => {
        H(), G(K);
      }, z = (K) => {
        H(), W(K);
      };
      if (this.pendingMcpResponses.set(J, { resolve: B, reject: z }), X.onmessage) X.onmessage($.message);
      else {
        H(), W(Error("No message handler registered"));
        return;
      }
    });
  }
};
var yU = xU(TU);
var a9 = Buffer.from('{"type":"attribution-snapshot"');
var iU = Buffer.from('{"type":"system"');
var f4 = 10;
var nU = Buffer.from([f4]);
var i1 = {};
uQ(i1, { void: () => BF, util: () => d, unknown: () => WF, union: () => VF, undefined: () => YF, tuple: () => LF, transformer: () => IF, symbol: () => XF, string: () => EJ, strictObject: () => KF, setErrorMap: () => bL, set: () => OF, record: () => FF, quotelessJson: () => EL, promise: () => RF, preprocess: () => bF, pipeline: () => ZF, ostring: () => CF, optional: () => EF, onumber: () => SF, oboolean: () => _F, objectUtil: () => W$, object: () => z$, number: () => PJ, nullable: () => PF, null: () => JF, never: () => HF, nativeEnum: () => jF, nan: () => eL, map: () => NF, makeIssue: () => l4, literal: () => MF, lazy: () => wF, late: () => aL, isValid: () => m1, isDirty: () => J8, isAsync: () => t6, isAborted: () => Y8, intersection: () => UF, instanceof: () => sL, getParsedType: () => w1, getErrorMap: () => r6, function: () => DF, enum: () => AF, effect: () => IF, discriminatedUnion: () => qF, defaultErrorMap: () => Z1, datetimeRegex: () => jJ, date: () => $F, custom: () => IJ, coerce: () => kF, boolean: () => bJ, bigint: () => QF, array: () => zF, any: () => GF, addIssueToContext: () => E, ZodVoid: () => p4, ZodUnknown: () => l1, ZodUnion: () => X4, ZodUndefined: () => Q4, ZodType: () => p, ZodTuple: () => A1, ZodTransformer: () => G1, ZodSymbol: () => c4, ZodString: () => X1, ZodSet: () => D6, ZodSchema: () => p, ZodRecord: () => d4, ZodReadonly: () => z4, ZodPromise: () => w6, ZodPipeline: () => o4, ZodParsedType: () => I, ZodOptional: () => m0, ZodObject: () => U0, ZodNumber: () => c1, ZodNullable: () => S1, ZodNull: () => $4, ZodNever: () => M1, ZodNativeEnum: () => W4, ZodNaN: () => n4, ZodMap: () => i4, ZodLiteral: () => G4, ZodLazy: () => J4, ZodIssueCode: () => A, ZodIntersection: () => Y4, ZodFunction: () => s6, ZodFirstPartyTypeKind: () => j, ZodError: () => x0, ZodEnum: () => d1, ZodEffects: () => G1, ZodDiscriminatedUnion: () => G8, ZodDefault: () => H4, ZodDate: () => N6, ZodCatch: () => B4, ZodBranded: () => W8, ZodBoolean: () => e6, ZodBigInt: () => p1, ZodArray: () => Y1, ZodAny: () => O6, Schema: () => p, ParseStatus: () => A0, OK: () => b0, NEVER: () => vF, INVALID: () => x, EMPTY_PATH: () => ZL, DIRTY: () => F6, BRAND: () => tL });
var d;
(function(Q) {
  Q.assertEqual = (J) => {
  };
  function $(J) {
  }
  Q.assertIs = $;
  function X(J) {
    throw Error();
  }
  Q.assertNever = X, Q.arrayToEnum = (J) => {
    let G = {};
    for (let W of J) G[W] = W;
    return G;
  }, Q.getValidEnumValues = (J) => {
    let G = Q.objectKeys(J).filter((H) => typeof J[J[H]] !== "number"), W = {};
    for (let H of G) W[H] = J[H];
    return Q.objectValues(W);
  }, Q.objectValues = (J) => {
    return Q.objectKeys(J).map(function(G) {
      return J[G];
    });
  }, Q.objectKeys = typeof Object.keys === "function" ? (J) => Object.keys(J) : (J) => {
    let G = [];
    for (let W in J) if (Object.prototype.hasOwnProperty.call(J, W)) G.push(W);
    return G;
  }, Q.find = (J, G) => {
    for (let W of J) if (G(W)) return W;
    return;
  }, Q.isInteger = typeof Number.isInteger === "function" ? (J) => Number.isInteger(J) : (J) => typeof J === "number" && Number.isFinite(J) && Math.floor(J) === J;
  function Y(J, G = " | ") {
    return J.map((W) => typeof W === "string" ? `'${W}'` : W).join(G);
  }
  Q.joinValues = Y, Q.jsonStringifyReplacer = (J, G) => {
    if (typeof G === "bigint") return G.toString();
    return G;
  };
})(d || (d = {}));
var W$;
(function(Q) {
  Q.mergeShapes = ($, X) => {
    return { ...$, ...X };
  };
})(W$ || (W$ = {}));
var I = d.arrayToEnum(["string", "nan", "number", "integer", "float", "boolean", "date", "bigint", "symbol", "function", "undefined", "null", "array", "object", "unknown", "promise", "void", "never", "map", "set"]);
var w1 = (Q) => {
  switch (typeof Q) {
    case "undefined":
      return I.undefined;
    case "string":
      return I.string;
    case "number":
      return Number.isNaN(Q) ? I.nan : I.number;
    case "boolean":
      return I.boolean;
    case "function":
      return I.function;
    case "bigint":
      return I.bigint;
    case "symbol":
      return I.symbol;
    case "object":
      if (Array.isArray(Q)) return I.array;
      if (Q === null) return I.null;
      if (Q.then && typeof Q.then === "function" && Q.catch && typeof Q.catch === "function") return I.promise;
      if (typeof Map < "u" && Q instanceof Map) return I.map;
      if (typeof Set < "u" && Q instanceof Set) return I.set;
      if (typeof Date < "u" && Q instanceof Date) return I.date;
      return I.object;
    default:
      return I.unknown;
  }
};
var A = d.arrayToEnum(["invalid_type", "invalid_literal", "custom", "invalid_union", "invalid_union_discriminator", "invalid_enum_value", "unrecognized_keys", "invalid_arguments", "invalid_return_type", "invalid_date", "invalid_string", "too_small", "too_big", "invalid_intersection_types", "not_multiple_of", "not_finite"]);
var EL = (Q) => {
  return JSON.stringify(Q, null, 2).replace(/"([^"]+)":/g, "$1:");
};
var x0 = class _x0 extends Error {
  get errors() {
    return this.issues;
  }
  constructor(Q) {
    super();
    this.issues = [], this.addIssue = (X) => {
      this.issues = [...this.issues, X];
    }, this.addIssues = (X = []) => {
      this.issues = [...this.issues, ...X];
    };
    let $ = new.target.prototype;
    if (Object.setPrototypeOf) Object.setPrototypeOf(this, $);
    else this.__proto__ = $;
    this.name = "ZodError", this.issues = Q;
  }
  format(Q) {
    let $ = Q || function(J) {
      return J.message;
    }, X = { _errors: [] }, Y = (J) => {
      for (let G of J.issues) if (G.code === "invalid_union") G.unionErrors.map(Y);
      else if (G.code === "invalid_return_type") Y(G.returnTypeError);
      else if (G.code === "invalid_arguments") Y(G.argumentsError);
      else if (G.path.length === 0) X._errors.push($(G));
      else {
        let W = X, H = 0;
        while (H < G.path.length) {
          let B = G.path[H];
          if (H !== G.path.length - 1) W[B] = W[B] || { _errors: [] };
          else W[B] = W[B] || { _errors: [] }, W[B]._errors.push($(G));
          W = W[B], H++;
        }
      }
    };
    return Y(this), X;
  }
  static assert(Q) {
    if (!(Q instanceof _x0)) throw Error(`Not a ZodError: ${Q}`);
  }
  toString() {
    return this.message;
  }
  get message() {
    return JSON.stringify(this.issues, d.jsonStringifyReplacer, 2);
  }
  get isEmpty() {
    return this.issues.length === 0;
  }
  flatten(Q = ($) => $.message) {
    let $ = {}, X = [];
    for (let Y of this.issues) if (Y.path.length > 0) {
      let J = Y.path[0];
      $[J] = $[J] || [], $[J].push(Q(Y));
    } else X.push(Q(Y));
    return { formErrors: X, fieldErrors: $ };
  }
  get formErrors() {
    return this.flatten();
  }
};
x0.create = (Q) => {
  return new x0(Q);
};
var PL = (Q, $) => {
  let X;
  switch (Q.code) {
    case A.invalid_type:
      if (Q.received === I.undefined) X = "Required";
      else X = `Expected ${Q.expected}, received ${Q.received}`;
      break;
    case A.invalid_literal:
      X = `Invalid literal value, expected ${JSON.stringify(Q.expected, d.jsonStringifyReplacer)}`;
      break;
    case A.unrecognized_keys:
      X = `Unrecognized key(s) in object: ${d.joinValues(Q.keys, ", ")}`;
      break;
    case A.invalid_union:
      X = "Invalid input";
      break;
    case A.invalid_union_discriminator:
      X = `Invalid discriminator value. Expected ${d.joinValues(Q.options)}`;
      break;
    case A.invalid_enum_value:
      X = `Invalid enum value. Expected ${d.joinValues(Q.options)}, received '${Q.received}'`;
      break;
    case A.invalid_arguments:
      X = "Invalid function arguments";
      break;
    case A.invalid_return_type:
      X = "Invalid function return type";
      break;
    case A.invalid_date:
      X = "Invalid date";
      break;
    case A.invalid_string:
      if (typeof Q.validation === "object") if ("includes" in Q.validation) {
        if (X = `Invalid input: must include "${Q.validation.includes}"`, typeof Q.validation.position === "number") X = `${X} at one or more positions greater than or equal to ${Q.validation.position}`;
      } else if ("startsWith" in Q.validation) X = `Invalid input: must start with "${Q.validation.startsWith}"`;
      else if ("endsWith" in Q.validation) X = `Invalid input: must end with "${Q.validation.endsWith}"`;
      else d.assertNever(Q.validation);
      else if (Q.validation !== "regex") X = `Invalid ${Q.validation}`;
      else X = "Invalid";
      break;
    case A.too_small:
      if (Q.type === "array") X = `Array must contain ${Q.exact ? "exactly" : Q.inclusive ? "at least" : "more than"} ${Q.minimum} element(s)`;
      else if (Q.type === "string") X = `String must contain ${Q.exact ? "exactly" : Q.inclusive ? "at least" : "over"} ${Q.minimum} character(s)`;
      else if (Q.type === "number") X = `Number must be ${Q.exact ? "exactly equal to " : Q.inclusive ? "greater than or equal to " : "greater than "}${Q.minimum}`;
      else if (Q.type === "bigint") X = `Number must be ${Q.exact ? "exactly equal to " : Q.inclusive ? "greater than or equal to " : "greater than "}${Q.minimum}`;
      else if (Q.type === "date") X = `Date must be ${Q.exact ? "exactly equal to " : Q.inclusive ? "greater than or equal to " : "greater than "}${new Date(Number(Q.minimum))}`;
      else X = "Invalid input";
      break;
    case A.too_big:
      if (Q.type === "array") X = `Array must contain ${Q.exact ? "exactly" : Q.inclusive ? "at most" : "less than"} ${Q.maximum} element(s)`;
      else if (Q.type === "string") X = `String must contain ${Q.exact ? "exactly" : Q.inclusive ? "at most" : "under"} ${Q.maximum} character(s)`;
      else if (Q.type === "number") X = `Number must be ${Q.exact ? "exactly" : Q.inclusive ? "less than or equal to" : "less than"} ${Q.maximum}`;
      else if (Q.type === "bigint") X = `BigInt must be ${Q.exact ? "exactly" : Q.inclusive ? "less than or equal to" : "less than"} ${Q.maximum}`;
      else if (Q.type === "date") X = `Date must be ${Q.exact ? "exactly" : Q.inclusive ? "smaller than or equal to" : "smaller than"} ${new Date(Number(Q.maximum))}`;
      else X = "Invalid input";
      break;
    case A.custom:
      X = "Invalid input";
      break;
    case A.invalid_intersection_types:
      X = "Intersection results could not be merged";
      break;
    case A.not_multiple_of:
      X = `Number must be a multiple of ${Q.multipleOf}`;
      break;
    case A.not_finite:
      X = "Number must be finite";
      break;
    default:
      X = $.defaultError, d.assertNever(Q);
  }
  return { message: X };
};
var Z1 = PL;
var OJ = Z1;
function bL(Q) {
  OJ = Q;
}
function r6() {
  return OJ;
}
var l4 = (Q) => {
  let { data: $, path: X, errorMaps: Y, issueData: J } = Q, G = [...X, ...J.path || []], W = { ...J, path: G };
  if (J.message !== void 0) return { ...J, path: G, message: J.message };
  let H = "", B = Y.filter((z) => !!z).slice().reverse();
  for (let z of B) H = z(W, { data: $, defaultError: H }).message;
  return { ...J, path: G, message: H };
};
var ZL = [];
function E(Q, $) {
  let X = r6(), Y = l4({ issueData: $, data: Q.data, path: Q.path, errorMaps: [Q.common.contextualErrorMap, Q.schemaErrorMap, X, X === Z1 ? void 0 : Z1].filter((J) => !!J) });
  Q.common.issues.push(Y);
}
var A0 = class _A0 {
  constructor() {
    this.value = "valid";
  }
  dirty() {
    if (this.value === "valid") this.value = "dirty";
  }
  abort() {
    if (this.value !== "aborted") this.value = "aborted";
  }
  static mergeArray(Q, $) {
    let X = [];
    for (let Y of $) {
      if (Y.status === "aborted") return x;
      if (Y.status === "dirty") Q.dirty();
      X.push(Y.value);
    }
    return { status: Q.value, value: X };
  }
  static async mergeObjectAsync(Q, $) {
    let X = [];
    for (let Y of $) {
      let J = await Y.key, G = await Y.value;
      X.push({ key: J, value: G });
    }
    return _A0.mergeObjectSync(Q, X);
  }
  static mergeObjectSync(Q, $) {
    let X = {};
    for (let Y of $) {
      let { key: J, value: G } = Y;
      if (J.status === "aborted") return x;
      if (G.status === "aborted") return x;
      if (J.status === "dirty") Q.dirty();
      if (G.status === "dirty") Q.dirty();
      if (J.value !== "__proto__" && (typeof G.value < "u" || Y.alwaysSet)) X[J.value] = G.value;
    }
    return { status: Q.value, value: X };
  }
};
var x = Object.freeze({ status: "aborted" });
var F6 = (Q) => ({ status: "dirty", value: Q });
var b0 = (Q) => ({ status: "valid", value: Q });
var Y8 = (Q) => Q.status === "aborted";
var J8 = (Q) => Q.status === "dirty";
var m1 = (Q) => Q.status === "valid";
var t6 = (Q) => typeof Promise < "u" && Q instanceof Promise;
var C;
(function(Q) {
  Q.errToObj = ($) => typeof $ === "string" ? { message: $ } : $ || {}, Q.toString = ($) => typeof $ === "string" ? $ : $?.message;
})(C || (C = {}));
var J1 = class {
  constructor(Q, $, X, Y) {
    this._cachedPath = [], this.parent = Q, this.data = $, this._path = X, this._key = Y;
  }
  get path() {
    if (!this._cachedPath.length) if (Array.isArray(this._key)) this._cachedPath.push(...this._path, ...this._key);
    else this._cachedPath.push(...this._path, this._key);
    return this._cachedPath;
  }
};
var DJ = (Q, $) => {
  if (m1($)) return { success: true, data: $.value };
  else {
    if (!Q.common.issues.length) throw Error("Validation failed but no issues detected.");
    return { success: false, get error() {
      if (this._error) return this._error;
      let X = new x0(Q.common.issues);
      return this._error = X, this._error;
    } };
  }
};
function m(Q) {
  if (!Q) return {};
  let { errorMap: $, invalid_type_error: X, required_error: Y, description: J } = Q;
  if ($ && (X || Y)) throw Error(`Can't use "invalid_type_error" or "required_error" in conjunction with custom error map.`);
  if ($) return { errorMap: $, description: J };
  return { errorMap: (W, H) => {
    let { message: B } = Q;
    if (W.code === "invalid_enum_value") return { message: B ?? H.defaultError };
    if (typeof H.data > "u") return { message: B ?? Y ?? H.defaultError };
    if (W.code !== "invalid_type") return { message: H.defaultError };
    return { message: B ?? X ?? H.defaultError };
  }, description: J };
}
var p = class {
  get description() {
    return this._def.description;
  }
  _getType(Q) {
    return w1(Q.data);
  }
  _getOrReturnCtx(Q, $) {
    return $ || { common: Q.parent.common, data: Q.data, parsedType: w1(Q.data), schemaErrorMap: this._def.errorMap, path: Q.path, parent: Q.parent };
  }
  _processInputParams(Q) {
    return { status: new A0(), ctx: { common: Q.parent.common, data: Q.data, parsedType: w1(Q.data), schemaErrorMap: this._def.errorMap, path: Q.path, parent: Q.parent } };
  }
  _parseSync(Q) {
    let $ = this._parse(Q);
    if (t6($)) throw Error("Synchronous parse encountered promise.");
    return $;
  }
  _parseAsync(Q) {
    let $ = this._parse(Q);
    return Promise.resolve($);
  }
  parse(Q, $) {
    let X = this.safeParse(Q, $);
    if (X.success) return X.data;
    throw X.error;
  }
  safeParse(Q, $) {
    let X = { common: { issues: [], async: $?.async ?? false, contextualErrorMap: $?.errorMap }, path: $?.path || [], schemaErrorMap: this._def.errorMap, parent: null, data: Q, parsedType: w1(Q) }, Y = this._parseSync({ data: Q, path: X.path, parent: X });
    return DJ(X, Y);
  }
  "~validate"(Q) {
    let $ = { common: { issues: [], async: !!this["~standard"].async }, path: [], schemaErrorMap: this._def.errorMap, parent: null, data: Q, parsedType: w1(Q) };
    if (!this["~standard"].async) try {
      let X = this._parseSync({ data: Q, path: [], parent: $ });
      return m1(X) ? { value: X.value } : { issues: $.common.issues };
    } catch (X) {
      if (X?.message?.toLowerCase()?.includes("encountered")) this["~standard"].async = true;
      $.common = { issues: [], async: true };
    }
    return this._parseAsync({ data: Q, path: [], parent: $ }).then((X) => m1(X) ? { value: X.value } : { issues: $.common.issues });
  }
  async parseAsync(Q, $) {
    let X = await this.safeParseAsync(Q, $);
    if (X.success) return X.data;
    throw X.error;
  }
  async safeParseAsync(Q, $) {
    let X = { common: { issues: [], contextualErrorMap: $?.errorMap, async: true }, path: $?.path || [], schemaErrorMap: this._def.errorMap, parent: null, data: Q, parsedType: w1(Q) }, Y = this._parse({ data: Q, path: X.path, parent: X }), J = await (t6(Y) ? Y : Promise.resolve(Y));
    return DJ(X, J);
  }
  refine(Q, $) {
    let X = (Y) => {
      if (typeof $ === "string" || typeof $ > "u") return { message: $ };
      else if (typeof $ === "function") return $(Y);
      else return $;
    };
    return this._refinement((Y, J) => {
      let G = Q(Y), W = () => J.addIssue({ code: A.custom, ...X(Y) });
      if (typeof Promise < "u" && G instanceof Promise) return G.then((H) => {
        if (!H) return W(), false;
        else return true;
      });
      if (!G) return W(), false;
      else return true;
    });
  }
  refinement(Q, $) {
    return this._refinement((X, Y) => {
      if (!Q(X)) return Y.addIssue(typeof $ === "function" ? $(X, Y) : $), false;
      else return true;
    });
  }
  _refinement(Q) {
    return new G1({ schema: this, typeName: j.ZodEffects, effect: { type: "refinement", refinement: Q } });
  }
  superRefine(Q) {
    return this._refinement(Q);
  }
  constructor(Q) {
    this.spa = this.safeParseAsync, this._def = Q, this.parse = this.parse.bind(this), this.safeParse = this.safeParse.bind(this), this.parseAsync = this.parseAsync.bind(this), this.safeParseAsync = this.safeParseAsync.bind(this), this.spa = this.spa.bind(this), this.refine = this.refine.bind(this), this.refinement = this.refinement.bind(this), this.superRefine = this.superRefine.bind(this), this.optional = this.optional.bind(this), this.nullable = this.nullable.bind(this), this.nullish = this.nullish.bind(this), this.array = this.array.bind(this), this.promise = this.promise.bind(this), this.or = this.or.bind(this), this.and = this.and.bind(this), this.transform = this.transform.bind(this), this.brand = this.brand.bind(this), this.default = this.default.bind(this), this.catch = this.catch.bind(this), this.describe = this.describe.bind(this), this.pipe = this.pipe.bind(this), this.readonly = this.readonly.bind(this), this.isNullable = this.isNullable.bind(this), this.isOptional = this.isOptional.bind(this), this["~standard"] = { version: 1, vendor: "zod", validate: ($) => this["~validate"]($) };
  }
  optional() {
    return m0.create(this, this._def);
  }
  nullable() {
    return S1.create(this, this._def);
  }
  nullish() {
    return this.nullable().optional();
  }
  array() {
    return Y1.create(this);
  }
  promise() {
    return w6.create(this, this._def);
  }
  or(Q) {
    return X4.create([this, Q], this._def);
  }
  and(Q) {
    return Y4.create(this, Q, this._def);
  }
  transform(Q) {
    return new G1({ ...m(this._def), schema: this, typeName: j.ZodEffects, effect: { type: "transform", transform: Q } });
  }
  default(Q) {
    let $ = typeof Q === "function" ? Q : () => Q;
    return new H4({ ...m(this._def), innerType: this, defaultValue: $, typeName: j.ZodDefault });
  }
  brand() {
    return new W8({ typeName: j.ZodBranded, type: this, ...m(this._def) });
  }
  catch(Q) {
    let $ = typeof Q === "function" ? Q : () => Q;
    return new B4({ ...m(this._def), innerType: this, catchValue: $, typeName: j.ZodCatch });
  }
  describe(Q) {
    return new this.constructor({ ...this._def, description: Q });
  }
  pipe(Q) {
    return o4.create(this, Q);
  }
  readonly() {
    return z4.create(this);
  }
  isOptional() {
    return this.safeParse(void 0).success;
  }
  isNullable() {
    return this.safeParse(null).success;
  }
};
var CL = /^c[^\s-]{8,}$/i;
var SL = /^[0-9a-z]+$/;
var _L = /^[0-9A-HJKMNP-TV-Z]{26}$/i;
var kL = /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/i;
var vL = /^[a-z0-9_-]{21}$/i;
var TL = /^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$/;
var xL = /^[-+]?P(?!$)(?:(?:[-+]?\d+Y)|(?:[-+]?\d+[.,]\d+Y$))?(?:(?:[-+]?\d+M)|(?:[-+]?\d+[.,]\d+M$))?(?:(?:[-+]?\d+W)|(?:[-+]?\d+[.,]\d+W$))?(?:(?:[-+]?\d+D)|(?:[-+]?\d+[.,]\d+D$))?(?:T(?=[\d+-])(?:(?:[-+]?\d+H)|(?:[-+]?\d+[.,]\d+H$))?(?:(?:[-+]?\d+M)|(?:[-+]?\d+[.,]\d+M$))?(?:[-+]?\d+(?:[.,]\d+)?S)?)??$/;
var yL = /^(?!\.)(?!.*\.\.)([A-Z0-9_'+\-\.]*)[A-Z0-9_+-]@([A-Z0-9][A-Z0-9\-]*\.)+[A-Z]{2,}$/i;
var gL = "^(\\p{Extended_Pictographic}|\\p{Emoji_Component})+$";
var H$;
var hL = /^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])$/;
var fL = /^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\/(3[0-2]|[12]?[0-9])$/;
var uL = /^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$/;
var mL = /^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))\/(12[0-8]|1[01][0-9]|[1-9]?[0-9])$/;
var lL = /^([0-9a-zA-Z+/]{4})*(([0-9a-zA-Z+/]{2}==)|([0-9a-zA-Z+/]{3}=))?$/;
var cL = /^([0-9a-zA-Z-_]{4})*(([0-9a-zA-Z-_]{2}(==)?)|([0-9a-zA-Z-_]{3}(=)?))?$/;
var MJ = "((\\d\\d[2468][048]|\\d\\d[13579][26]|\\d\\d0[48]|[02468][048]00|[13579][26]00)-02-29|\\d{4}-((0[13578]|1[02])-(0[1-9]|[12]\\d|3[01])|(0[469]|11)-(0[1-9]|[12]\\d|30)|(02)-(0[1-9]|1\\d|2[0-8])))";
var pL = new RegExp(`^${MJ}$`);
function AJ(Q) {
  let $ = "[0-5]\\d";
  if (Q.precision) $ = `${$}\\.\\d{${Q.precision}}`;
  else if (Q.precision == null) $ = `${$}(\\.\\d+)?`;
  let X = Q.precision ? "+" : "?";
  return `([01]\\d|2[0-3]):[0-5]\\d(:${$})${X}`;
}
function dL(Q) {
  return new RegExp(`^${AJ(Q)}$`);
}
function jJ(Q) {
  let $ = `${MJ}T${AJ(Q)}`, X = [];
  if (X.push(Q.local ? "Z?" : "Z"), Q.offset) X.push("([+-]\\d{2}:?\\d{2})");
  return $ = `${$}(${X.join("|")})`, new RegExp(`^${$}$`);
}
function iL(Q, $) {
  if (($ === "v4" || !$) && hL.test(Q)) return true;
  if (($ === "v6" || !$) && uL.test(Q)) return true;
  return false;
}
function nL(Q, $) {
  if (!TL.test(Q)) return false;
  try {
    let [X] = Q.split(".");
    if (!X) return false;
    let Y = X.replace(/-/g, "+").replace(/_/g, "/").padEnd(X.length + (4 - X.length % 4) % 4, "="), J = JSON.parse(atob(Y));
    if (typeof J !== "object" || J === null) return false;
    if ("typ" in J && J?.typ !== "JWT") return false;
    if (!J.alg) return false;
    if ($ && J.alg !== $) return false;
    return true;
  } catch {
    return false;
  }
}
function oL(Q, $) {
  if (($ === "v4" || !$) && fL.test(Q)) return true;
  if (($ === "v6" || !$) && mL.test(Q)) return true;
  return false;
}
var X1 = class _X1 extends p {
  _parse(Q) {
    if (this._def.coerce) Q.data = String(Q.data);
    if (this._getType(Q) !== I.string) {
      let J = this._getOrReturnCtx(Q);
      return E(J, { code: A.invalid_type, expected: I.string, received: J.parsedType }), x;
    }
    let X = new A0(), Y = void 0;
    for (let J of this._def.checks) if (J.kind === "min") {
      if (Q.data.length < J.value) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.too_small, minimum: J.value, type: "string", inclusive: true, exact: false, message: J.message }), X.dirty();
    } else if (J.kind === "max") {
      if (Q.data.length > J.value) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.too_big, maximum: J.value, type: "string", inclusive: true, exact: false, message: J.message }), X.dirty();
    } else if (J.kind === "length") {
      let G = Q.data.length > J.value, W = Q.data.length < J.value;
      if (G || W) {
        if (Y = this._getOrReturnCtx(Q, Y), G) E(Y, { code: A.too_big, maximum: J.value, type: "string", inclusive: true, exact: true, message: J.message });
        else if (W) E(Y, { code: A.too_small, minimum: J.value, type: "string", inclusive: true, exact: true, message: J.message });
        X.dirty();
      }
    } else if (J.kind === "email") {
      if (!yL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "email", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "emoji") {
      if (!H$) H$ = new RegExp(gL, "u");
      if (!H$.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "emoji", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "uuid") {
      if (!kL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "uuid", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "nanoid") {
      if (!vL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "nanoid", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "cuid") {
      if (!CL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "cuid", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "cuid2") {
      if (!SL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "cuid2", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "ulid") {
      if (!_L.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "ulid", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "url") try {
      new URL(Q.data);
    } catch {
      Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "url", code: A.invalid_string, message: J.message }), X.dirty();
    }
    else if (J.kind === "regex") {
      if (J.regex.lastIndex = 0, !J.regex.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "regex", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "trim") Q.data = Q.data.trim();
    else if (J.kind === "includes") {
      if (!Q.data.includes(J.value, J.position)) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.invalid_string, validation: { includes: J.value, position: J.position }, message: J.message }), X.dirty();
    } else if (J.kind === "toLowerCase") Q.data = Q.data.toLowerCase();
    else if (J.kind === "toUpperCase") Q.data = Q.data.toUpperCase();
    else if (J.kind === "startsWith") {
      if (!Q.data.startsWith(J.value)) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.invalid_string, validation: { startsWith: J.value }, message: J.message }), X.dirty();
    } else if (J.kind === "endsWith") {
      if (!Q.data.endsWith(J.value)) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.invalid_string, validation: { endsWith: J.value }, message: J.message }), X.dirty();
    } else if (J.kind === "datetime") {
      if (!jJ(J).test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.invalid_string, validation: "datetime", message: J.message }), X.dirty();
    } else if (J.kind === "date") {
      if (!pL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.invalid_string, validation: "date", message: J.message }), X.dirty();
    } else if (J.kind === "time") {
      if (!dL(J).test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.invalid_string, validation: "time", message: J.message }), X.dirty();
    } else if (J.kind === "duration") {
      if (!xL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "duration", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "ip") {
      if (!iL(Q.data, J.version)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "ip", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "jwt") {
      if (!nL(Q.data, J.alg)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "jwt", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "cidr") {
      if (!oL(Q.data, J.version)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "cidr", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "base64") {
      if (!lL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "base64", code: A.invalid_string, message: J.message }), X.dirty();
    } else if (J.kind === "base64url") {
      if (!cL.test(Q.data)) Y = this._getOrReturnCtx(Q, Y), E(Y, { validation: "base64url", code: A.invalid_string, message: J.message }), X.dirty();
    } else d.assertNever(J);
    return { status: X.value, value: Q.data };
  }
  _regex(Q, $, X) {
    return this.refinement((Y) => Q.test(Y), { validation: $, code: A.invalid_string, ...C.errToObj(X) });
  }
  _addCheck(Q) {
    return new _X1({ ...this._def, checks: [...this._def.checks, Q] });
  }
  email(Q) {
    return this._addCheck({ kind: "email", ...C.errToObj(Q) });
  }
  url(Q) {
    return this._addCheck({ kind: "url", ...C.errToObj(Q) });
  }
  emoji(Q) {
    return this._addCheck({ kind: "emoji", ...C.errToObj(Q) });
  }
  uuid(Q) {
    return this._addCheck({ kind: "uuid", ...C.errToObj(Q) });
  }
  nanoid(Q) {
    return this._addCheck({ kind: "nanoid", ...C.errToObj(Q) });
  }
  cuid(Q) {
    return this._addCheck({ kind: "cuid", ...C.errToObj(Q) });
  }
  cuid2(Q) {
    return this._addCheck({ kind: "cuid2", ...C.errToObj(Q) });
  }
  ulid(Q) {
    return this._addCheck({ kind: "ulid", ...C.errToObj(Q) });
  }
  base64(Q) {
    return this._addCheck({ kind: "base64", ...C.errToObj(Q) });
  }
  base64url(Q) {
    return this._addCheck({ kind: "base64url", ...C.errToObj(Q) });
  }
  jwt(Q) {
    return this._addCheck({ kind: "jwt", ...C.errToObj(Q) });
  }
  ip(Q) {
    return this._addCheck({ kind: "ip", ...C.errToObj(Q) });
  }
  cidr(Q) {
    return this._addCheck({ kind: "cidr", ...C.errToObj(Q) });
  }
  datetime(Q) {
    if (typeof Q === "string") return this._addCheck({ kind: "datetime", precision: null, offset: false, local: false, message: Q });
    return this._addCheck({ kind: "datetime", precision: typeof Q?.precision > "u" ? null : Q?.precision, offset: Q?.offset ?? false, local: Q?.local ?? false, ...C.errToObj(Q?.message) });
  }
  date(Q) {
    return this._addCheck({ kind: "date", message: Q });
  }
  time(Q) {
    if (typeof Q === "string") return this._addCheck({ kind: "time", precision: null, message: Q });
    return this._addCheck({ kind: "time", precision: typeof Q?.precision > "u" ? null : Q?.precision, ...C.errToObj(Q?.message) });
  }
  duration(Q) {
    return this._addCheck({ kind: "duration", ...C.errToObj(Q) });
  }
  regex(Q, $) {
    return this._addCheck({ kind: "regex", regex: Q, ...C.errToObj($) });
  }
  includes(Q, $) {
    return this._addCheck({ kind: "includes", value: Q, position: $?.position, ...C.errToObj($?.message) });
  }
  startsWith(Q, $) {
    return this._addCheck({ kind: "startsWith", value: Q, ...C.errToObj($) });
  }
  endsWith(Q, $) {
    return this._addCheck({ kind: "endsWith", value: Q, ...C.errToObj($) });
  }
  min(Q, $) {
    return this._addCheck({ kind: "min", value: Q, ...C.errToObj($) });
  }
  max(Q, $) {
    return this._addCheck({ kind: "max", value: Q, ...C.errToObj($) });
  }
  length(Q, $) {
    return this._addCheck({ kind: "length", value: Q, ...C.errToObj($) });
  }
  nonempty(Q) {
    return this.min(1, C.errToObj(Q));
  }
  trim() {
    return new _X1({ ...this._def, checks: [...this._def.checks, { kind: "trim" }] });
  }
  toLowerCase() {
    return new _X1({ ...this._def, checks: [...this._def.checks, { kind: "toLowerCase" }] });
  }
  toUpperCase() {
    return new _X1({ ...this._def, checks: [...this._def.checks, { kind: "toUpperCase" }] });
  }
  get isDatetime() {
    return !!this._def.checks.find((Q) => Q.kind === "datetime");
  }
  get isDate() {
    return !!this._def.checks.find((Q) => Q.kind === "date");
  }
  get isTime() {
    return !!this._def.checks.find((Q) => Q.kind === "time");
  }
  get isDuration() {
    return !!this._def.checks.find((Q) => Q.kind === "duration");
  }
  get isEmail() {
    return !!this._def.checks.find((Q) => Q.kind === "email");
  }
  get isURL() {
    return !!this._def.checks.find((Q) => Q.kind === "url");
  }
  get isEmoji() {
    return !!this._def.checks.find((Q) => Q.kind === "emoji");
  }
  get isUUID() {
    return !!this._def.checks.find((Q) => Q.kind === "uuid");
  }
  get isNANOID() {
    return !!this._def.checks.find((Q) => Q.kind === "nanoid");
  }
  get isCUID() {
    return !!this._def.checks.find((Q) => Q.kind === "cuid");
  }
  get isCUID2() {
    return !!this._def.checks.find((Q) => Q.kind === "cuid2");
  }
  get isULID() {
    return !!this._def.checks.find((Q) => Q.kind === "ulid");
  }
  get isIP() {
    return !!this._def.checks.find((Q) => Q.kind === "ip");
  }
  get isCIDR() {
    return !!this._def.checks.find((Q) => Q.kind === "cidr");
  }
  get isBase64() {
    return !!this._def.checks.find((Q) => Q.kind === "base64");
  }
  get isBase64url() {
    return !!this._def.checks.find((Q) => Q.kind === "base64url");
  }
  get minLength() {
    let Q = null;
    for (let $ of this._def.checks) if ($.kind === "min") {
      if (Q === null || $.value > Q) Q = $.value;
    }
    return Q;
  }
  get maxLength() {
    let Q = null;
    for (let $ of this._def.checks) if ($.kind === "max") {
      if (Q === null || $.value < Q) Q = $.value;
    }
    return Q;
  }
};
X1.create = (Q) => {
  return new X1({ checks: [], typeName: j.ZodString, coerce: Q?.coerce ?? false, ...m(Q) });
};
function rL(Q, $) {
  let X = (Q.toString().split(".")[1] || "").length, Y = ($.toString().split(".")[1] || "").length, J = X > Y ? X : Y, G = Number.parseInt(Q.toFixed(J).replace(".", "")), W = Number.parseInt($.toFixed(J).replace(".", ""));
  return G % W / 10 ** J;
}
var c1 = class _c1 extends p {
  constructor() {
    super(...arguments);
    this.min = this.gte, this.max = this.lte, this.step = this.multipleOf;
  }
  _parse(Q) {
    if (this._def.coerce) Q.data = Number(Q.data);
    if (this._getType(Q) !== I.number) {
      let J = this._getOrReturnCtx(Q);
      return E(J, { code: A.invalid_type, expected: I.number, received: J.parsedType }), x;
    }
    let X = void 0, Y = new A0();
    for (let J of this._def.checks) if (J.kind === "int") {
      if (!d.isInteger(Q.data)) X = this._getOrReturnCtx(Q, X), E(X, { code: A.invalid_type, expected: "integer", received: "float", message: J.message }), Y.dirty();
    } else if (J.kind === "min") {
      if (J.inclusive ? Q.data < J.value : Q.data <= J.value) X = this._getOrReturnCtx(Q, X), E(X, { code: A.too_small, minimum: J.value, type: "number", inclusive: J.inclusive, exact: false, message: J.message }), Y.dirty();
    } else if (J.kind === "max") {
      if (J.inclusive ? Q.data > J.value : Q.data >= J.value) X = this._getOrReturnCtx(Q, X), E(X, { code: A.too_big, maximum: J.value, type: "number", inclusive: J.inclusive, exact: false, message: J.message }), Y.dirty();
    } else if (J.kind === "multipleOf") {
      if (rL(Q.data, J.value) !== 0) X = this._getOrReturnCtx(Q, X), E(X, { code: A.not_multiple_of, multipleOf: J.value, message: J.message }), Y.dirty();
    } else if (J.kind === "finite") {
      if (!Number.isFinite(Q.data)) X = this._getOrReturnCtx(Q, X), E(X, { code: A.not_finite, message: J.message }), Y.dirty();
    } else d.assertNever(J);
    return { status: Y.value, value: Q.data };
  }
  gte(Q, $) {
    return this.setLimit("min", Q, true, C.toString($));
  }
  gt(Q, $) {
    return this.setLimit("min", Q, false, C.toString($));
  }
  lte(Q, $) {
    return this.setLimit("max", Q, true, C.toString($));
  }
  lt(Q, $) {
    return this.setLimit("max", Q, false, C.toString($));
  }
  setLimit(Q, $, X, Y) {
    return new _c1({ ...this._def, checks: [...this._def.checks, { kind: Q, value: $, inclusive: X, message: C.toString(Y) }] });
  }
  _addCheck(Q) {
    return new _c1({ ...this._def, checks: [...this._def.checks, Q] });
  }
  int(Q) {
    return this._addCheck({ kind: "int", message: C.toString(Q) });
  }
  positive(Q) {
    return this._addCheck({ kind: "min", value: 0, inclusive: false, message: C.toString(Q) });
  }
  negative(Q) {
    return this._addCheck({ kind: "max", value: 0, inclusive: false, message: C.toString(Q) });
  }
  nonpositive(Q) {
    return this._addCheck({ kind: "max", value: 0, inclusive: true, message: C.toString(Q) });
  }
  nonnegative(Q) {
    return this._addCheck({ kind: "min", value: 0, inclusive: true, message: C.toString(Q) });
  }
  multipleOf(Q, $) {
    return this._addCheck({ kind: "multipleOf", value: Q, message: C.toString($) });
  }
  finite(Q) {
    return this._addCheck({ kind: "finite", message: C.toString(Q) });
  }
  safe(Q) {
    return this._addCheck({ kind: "min", inclusive: true, value: Number.MIN_SAFE_INTEGER, message: C.toString(Q) })._addCheck({ kind: "max", inclusive: true, value: Number.MAX_SAFE_INTEGER, message: C.toString(Q) });
  }
  get minValue() {
    let Q = null;
    for (let $ of this._def.checks) if ($.kind === "min") {
      if (Q === null || $.value > Q) Q = $.value;
    }
    return Q;
  }
  get maxValue() {
    let Q = null;
    for (let $ of this._def.checks) if ($.kind === "max") {
      if (Q === null || $.value < Q) Q = $.value;
    }
    return Q;
  }
  get isInt() {
    return !!this._def.checks.find((Q) => Q.kind === "int" || Q.kind === "multipleOf" && d.isInteger(Q.value));
  }
  get isFinite() {
    let Q = null, $ = null;
    for (let X of this._def.checks) if (X.kind === "finite" || X.kind === "int" || X.kind === "multipleOf") return true;
    else if (X.kind === "min") {
      if ($ === null || X.value > $) $ = X.value;
    } else if (X.kind === "max") {
      if (Q === null || X.value < Q) Q = X.value;
    }
    return Number.isFinite($) && Number.isFinite(Q);
  }
};
c1.create = (Q) => {
  return new c1({ checks: [], typeName: j.ZodNumber, coerce: Q?.coerce || false, ...m(Q) });
};
var p1 = class _p1 extends p {
  constructor() {
    super(...arguments);
    this.min = this.gte, this.max = this.lte;
  }
  _parse(Q) {
    if (this._def.coerce) try {
      Q.data = BigInt(Q.data);
    } catch {
      return this._getInvalidInput(Q);
    }
    if (this._getType(Q) !== I.bigint) return this._getInvalidInput(Q);
    let X = void 0, Y = new A0();
    for (let J of this._def.checks) if (J.kind === "min") {
      if (J.inclusive ? Q.data < J.value : Q.data <= J.value) X = this._getOrReturnCtx(Q, X), E(X, { code: A.too_small, type: "bigint", minimum: J.value, inclusive: J.inclusive, message: J.message }), Y.dirty();
    } else if (J.kind === "max") {
      if (J.inclusive ? Q.data > J.value : Q.data >= J.value) X = this._getOrReturnCtx(Q, X), E(X, { code: A.too_big, type: "bigint", maximum: J.value, inclusive: J.inclusive, message: J.message }), Y.dirty();
    } else if (J.kind === "multipleOf") {
      if (Q.data % J.value !== BigInt(0)) X = this._getOrReturnCtx(Q, X), E(X, { code: A.not_multiple_of, multipleOf: J.value, message: J.message }), Y.dirty();
    } else d.assertNever(J);
    return { status: Y.value, value: Q.data };
  }
  _getInvalidInput(Q) {
    let $ = this._getOrReturnCtx(Q);
    return E($, { code: A.invalid_type, expected: I.bigint, received: $.parsedType }), x;
  }
  gte(Q, $) {
    return this.setLimit("min", Q, true, C.toString($));
  }
  gt(Q, $) {
    return this.setLimit("min", Q, false, C.toString($));
  }
  lte(Q, $) {
    return this.setLimit("max", Q, true, C.toString($));
  }
  lt(Q, $) {
    return this.setLimit("max", Q, false, C.toString($));
  }
  setLimit(Q, $, X, Y) {
    return new _p1({ ...this._def, checks: [...this._def.checks, { kind: Q, value: $, inclusive: X, message: C.toString(Y) }] });
  }
  _addCheck(Q) {
    return new _p1({ ...this._def, checks: [...this._def.checks, Q] });
  }
  positive(Q) {
    return this._addCheck({ kind: "min", value: BigInt(0), inclusive: false, message: C.toString(Q) });
  }
  negative(Q) {
    return this._addCheck({ kind: "max", value: BigInt(0), inclusive: false, message: C.toString(Q) });
  }
  nonpositive(Q) {
    return this._addCheck({ kind: "max", value: BigInt(0), inclusive: true, message: C.toString(Q) });
  }
  nonnegative(Q) {
    return this._addCheck({ kind: "min", value: BigInt(0), inclusive: true, message: C.toString(Q) });
  }
  multipleOf(Q, $) {
    return this._addCheck({ kind: "multipleOf", value: Q, message: C.toString($) });
  }
  get minValue() {
    let Q = null;
    for (let $ of this._def.checks) if ($.kind === "min") {
      if (Q === null || $.value > Q) Q = $.value;
    }
    return Q;
  }
  get maxValue() {
    let Q = null;
    for (let $ of this._def.checks) if ($.kind === "max") {
      if (Q === null || $.value < Q) Q = $.value;
    }
    return Q;
  }
};
p1.create = (Q) => {
  return new p1({ checks: [], typeName: j.ZodBigInt, coerce: Q?.coerce ?? false, ...m(Q) });
};
var e6 = class extends p {
  _parse(Q) {
    if (this._def.coerce) Q.data = Boolean(Q.data);
    if (this._getType(Q) !== I.boolean) {
      let X = this._getOrReturnCtx(Q);
      return E(X, { code: A.invalid_type, expected: I.boolean, received: X.parsedType }), x;
    }
    return b0(Q.data);
  }
};
e6.create = (Q) => {
  return new e6({ typeName: j.ZodBoolean, coerce: Q?.coerce || false, ...m(Q) });
};
var N6 = class _N6 extends p {
  _parse(Q) {
    if (this._def.coerce) Q.data = new Date(Q.data);
    if (this._getType(Q) !== I.date) {
      let J = this._getOrReturnCtx(Q);
      return E(J, { code: A.invalid_type, expected: I.date, received: J.parsedType }), x;
    }
    if (Number.isNaN(Q.data.getTime())) {
      let J = this._getOrReturnCtx(Q);
      return E(J, { code: A.invalid_date }), x;
    }
    let X = new A0(), Y = void 0;
    for (let J of this._def.checks) if (J.kind === "min") {
      if (Q.data.getTime() < J.value) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.too_small, message: J.message, inclusive: true, exact: false, minimum: J.value, type: "date" }), X.dirty();
    } else if (J.kind === "max") {
      if (Q.data.getTime() > J.value) Y = this._getOrReturnCtx(Q, Y), E(Y, { code: A.too_big, message: J.message, inclusive: true, exact: false, maximum: J.value, type: "date" }), X.dirty();
    } else d.assertNever(J);
    return { status: X.value, value: new Date(Q.data.getTime()) };
  }
  _addCheck(Q) {
    return new _N6({ ...this._def, checks: [...this._def.checks, Q] });
  }
  min(Q, $) {
    return this._addCheck({ kind: "min", value: Q.getTime(), message: C.toString($) });
  }
  max(Q, $) {
    return this._addCheck({ kind: "max", value: Q.getTime(), message: C.toString($) });
  }
  get minDate() {
    let Q = null;
    for (let $ of this._def.checks) if ($.kind === "min") {
      if (Q === null || $.value > Q) Q = $.value;
    }
    return Q != null ? new Date(Q) : null;
  }
  get maxDate() {
    let Q = null;
    for (let $ of this._def.checks) if ($.kind === "max") {
      if (Q === null || $.value < Q) Q = $.value;
    }
    return Q != null ? new Date(Q) : null;
  }
};
N6.create = (Q) => {
  return new N6({ checks: [], coerce: Q?.coerce || false, typeName: j.ZodDate, ...m(Q) });
};
var c4 = class extends p {
  _parse(Q) {
    if (this._getType(Q) !== I.symbol) {
      let X = this._getOrReturnCtx(Q);
      return E(X, { code: A.invalid_type, expected: I.symbol, received: X.parsedType }), x;
    }
    return b0(Q.data);
  }
};
c4.create = (Q) => {
  return new c4({ typeName: j.ZodSymbol, ...m(Q) });
};
var Q4 = class extends p {
  _parse(Q) {
    if (this._getType(Q) !== I.undefined) {
      let X = this._getOrReturnCtx(Q);
      return E(X, { code: A.invalid_type, expected: I.undefined, received: X.parsedType }), x;
    }
    return b0(Q.data);
  }
};
Q4.create = (Q) => {
  return new Q4({ typeName: j.ZodUndefined, ...m(Q) });
};
var $4 = class extends p {
  _parse(Q) {
    if (this._getType(Q) !== I.null) {
      let X = this._getOrReturnCtx(Q);
      return E(X, { code: A.invalid_type, expected: I.null, received: X.parsedType }), x;
    }
    return b0(Q.data);
  }
};
$4.create = (Q) => {
  return new $4({ typeName: j.ZodNull, ...m(Q) });
};
var O6 = class extends p {
  constructor() {
    super(...arguments);
    this._any = true;
  }
  _parse(Q) {
    return b0(Q.data);
  }
};
O6.create = (Q) => {
  return new O6({ typeName: j.ZodAny, ...m(Q) });
};
var l1 = class extends p {
  constructor() {
    super(...arguments);
    this._unknown = true;
  }
  _parse(Q) {
    return b0(Q.data);
  }
};
l1.create = (Q) => {
  return new l1({ typeName: j.ZodUnknown, ...m(Q) });
};
var M1 = class extends p {
  _parse(Q) {
    let $ = this._getOrReturnCtx(Q);
    return E($, { code: A.invalid_type, expected: I.never, received: $.parsedType }), x;
  }
};
M1.create = (Q) => {
  return new M1({ typeName: j.ZodNever, ...m(Q) });
};
var p4 = class extends p {
  _parse(Q) {
    if (this._getType(Q) !== I.undefined) {
      let X = this._getOrReturnCtx(Q);
      return E(X, { code: A.invalid_type, expected: I.void, received: X.parsedType }), x;
    }
    return b0(Q.data);
  }
};
p4.create = (Q) => {
  return new p4({ typeName: j.ZodVoid, ...m(Q) });
};
var Y1 = class _Y1 extends p {
  _parse(Q) {
    let { ctx: $, status: X } = this._processInputParams(Q), Y = this._def;
    if ($.parsedType !== I.array) return E($, { code: A.invalid_type, expected: I.array, received: $.parsedType }), x;
    if (Y.exactLength !== null) {
      let G = $.data.length > Y.exactLength.value, W = $.data.length < Y.exactLength.value;
      if (G || W) E($, { code: G ? A.too_big : A.too_small, minimum: W ? Y.exactLength.value : void 0, maximum: G ? Y.exactLength.value : void 0, type: "array", inclusive: true, exact: true, message: Y.exactLength.message }), X.dirty();
    }
    if (Y.minLength !== null) {
      if ($.data.length < Y.minLength.value) E($, { code: A.too_small, minimum: Y.minLength.value, type: "array", inclusive: true, exact: false, message: Y.minLength.message }), X.dirty();
    }
    if (Y.maxLength !== null) {
      if ($.data.length > Y.maxLength.value) E($, { code: A.too_big, maximum: Y.maxLength.value, type: "array", inclusive: true, exact: false, message: Y.maxLength.message }), X.dirty();
    }
    if ($.common.async) return Promise.all([...$.data].map((G, W) => {
      return Y.type._parseAsync(new J1($, G, $.path, W));
    })).then((G) => {
      return A0.mergeArray(X, G);
    });
    let J = [...$.data].map((G, W) => {
      return Y.type._parseSync(new J1($, G, $.path, W));
    });
    return A0.mergeArray(X, J);
  }
  get element() {
    return this._def.type;
  }
  min(Q, $) {
    return new _Y1({ ...this._def, minLength: { value: Q, message: C.toString($) } });
  }
  max(Q, $) {
    return new _Y1({ ...this._def, maxLength: { value: Q, message: C.toString($) } });
  }
  length(Q, $) {
    return new _Y1({ ...this._def, exactLength: { value: Q, message: C.toString($) } });
  }
  nonempty(Q) {
    return this.min(1, Q);
  }
};
Y1.create = (Q, $) => {
  return new Y1({ type: Q, minLength: null, maxLength: null, exactLength: null, typeName: j.ZodArray, ...m($) });
};
function a6(Q) {
  if (Q instanceof U0) {
    let $ = {};
    for (let X in Q.shape) {
      let Y = Q.shape[X];
      $[X] = m0.create(a6(Y));
    }
    return new U0({ ...Q._def, shape: () => $ });
  } else if (Q instanceof Y1) return new Y1({ ...Q._def, type: a6(Q.element) });
  else if (Q instanceof m0) return m0.create(a6(Q.unwrap()));
  else if (Q instanceof S1) return S1.create(a6(Q.unwrap()));
  else if (Q instanceof A1) return A1.create(Q.items.map(($) => a6($)));
  else return Q;
}
var U0 = class _U0 extends p {
  constructor() {
    super(...arguments);
    this._cached = null, this.nonstrict = this.passthrough, this.augment = this.extend;
  }
  _getCached() {
    if (this._cached !== null) return this._cached;
    let Q = this._def.shape(), $ = d.objectKeys(Q);
    return this._cached = { shape: Q, keys: $ }, this._cached;
  }
  _parse(Q) {
    if (this._getType(Q) !== I.object) {
      let B = this._getOrReturnCtx(Q);
      return E(B, { code: A.invalid_type, expected: I.object, received: B.parsedType }), x;
    }
    let { status: X, ctx: Y } = this._processInputParams(Q), { shape: J, keys: G } = this._getCached(), W = [];
    if (!(this._def.catchall instanceof M1 && this._def.unknownKeys === "strip")) {
      for (let B in Y.data) if (!G.includes(B)) W.push(B);
    }
    let H = [];
    for (let B of G) {
      let z = J[B], K = Y.data[B];
      H.push({ key: { status: "valid", value: B }, value: z._parse(new J1(Y, K, Y.path, B)), alwaysSet: B in Y.data });
    }
    if (this._def.catchall instanceof M1) {
      let B = this._def.unknownKeys;
      if (B === "passthrough") for (let z of W) H.push({ key: { status: "valid", value: z }, value: { status: "valid", value: Y.data[z] } });
      else if (B === "strict") {
        if (W.length > 0) E(Y, { code: A.unrecognized_keys, keys: W }), X.dirty();
      } else if (B === "strip") ;
      else throw Error("Internal ZodObject error: invalid unknownKeys value.");
    } else {
      let B = this._def.catchall;
      for (let z of W) {
        let K = Y.data[z];
        H.push({ key: { status: "valid", value: z }, value: B._parse(new J1(Y, K, Y.path, z)), alwaysSet: z in Y.data });
      }
    }
    if (Y.common.async) return Promise.resolve().then(async () => {
      let B = [];
      for (let z of H) {
        let K = await z.key, U = await z.value;
        B.push({ key: K, value: U, alwaysSet: z.alwaysSet });
      }
      return B;
    }).then((B) => {
      return A0.mergeObjectSync(X, B);
    });
    else return A0.mergeObjectSync(X, H);
  }
  get shape() {
    return this._def.shape();
  }
  strict(Q) {
    return C.errToObj, new _U0({ ...this._def, unknownKeys: "strict", ...Q !== void 0 ? { errorMap: ($, X) => {
      let Y = this._def.errorMap?.($, X).message ?? X.defaultError;
      if ($.code === "unrecognized_keys") return { message: C.errToObj(Q).message ?? Y };
      return { message: Y };
    } } : {} });
  }
  strip() {
    return new _U0({ ...this._def, unknownKeys: "strip" });
  }
  passthrough() {
    return new _U0({ ...this._def, unknownKeys: "passthrough" });
  }
  extend(Q) {
    return new _U0({ ...this._def, shape: () => ({ ...this._def.shape(), ...Q }) });
  }
  merge(Q) {
    return new _U0({ unknownKeys: Q._def.unknownKeys, catchall: Q._def.catchall, shape: () => ({ ...this._def.shape(), ...Q._def.shape() }), typeName: j.ZodObject });
  }
  setKey(Q, $) {
    return this.augment({ [Q]: $ });
  }
  catchall(Q) {
    return new _U0({ ...this._def, catchall: Q });
  }
  pick(Q) {
    let $ = {};
    for (let X of d.objectKeys(Q)) if (Q[X] && this.shape[X]) $[X] = this.shape[X];
    return new _U0({ ...this._def, shape: () => $ });
  }
  omit(Q) {
    let $ = {};
    for (let X of d.objectKeys(this.shape)) if (!Q[X]) $[X] = this.shape[X];
    return new _U0({ ...this._def, shape: () => $ });
  }
  deepPartial() {
    return a6(this);
  }
  partial(Q) {
    let $ = {};
    for (let X of d.objectKeys(this.shape)) {
      let Y = this.shape[X];
      if (Q && !Q[X]) $[X] = Y;
      else $[X] = Y.optional();
    }
    return new _U0({ ...this._def, shape: () => $ });
  }
  required(Q) {
    let $ = {};
    for (let X of d.objectKeys(this.shape)) if (Q && !Q[X]) $[X] = this.shape[X];
    else {
      let J = this.shape[X];
      while (J instanceof m0) J = J._def.innerType;
      $[X] = J;
    }
    return new _U0({ ...this._def, shape: () => $ });
  }
  keyof() {
    return RJ(d.objectKeys(this.shape));
  }
};
U0.create = (Q, $) => {
  return new U0({ shape: () => Q, unknownKeys: "strip", catchall: M1.create(), typeName: j.ZodObject, ...m($) });
};
U0.strictCreate = (Q, $) => {
  return new U0({ shape: () => Q, unknownKeys: "strict", catchall: M1.create(), typeName: j.ZodObject, ...m($) });
};
U0.lazycreate = (Q, $) => {
  return new U0({ shape: Q, unknownKeys: "strip", catchall: M1.create(), typeName: j.ZodObject, ...m($) });
};
var X4 = class extends p {
  _parse(Q) {
    let { ctx: $ } = this._processInputParams(Q), X = this._def.options;
    function Y(J) {
      for (let W of J) if (W.result.status === "valid") return W.result;
      for (let W of J) if (W.result.status === "dirty") return $.common.issues.push(...W.ctx.common.issues), W.result;
      let G = J.map((W) => new x0(W.ctx.common.issues));
      return E($, { code: A.invalid_union, unionErrors: G }), x;
    }
    if ($.common.async) return Promise.all(X.map(async (J) => {
      let G = { ...$, common: { ...$.common, issues: [] }, parent: null };
      return { result: await J._parseAsync({ data: $.data, path: $.path, parent: G }), ctx: G };
    })).then(Y);
    else {
      let J = void 0, G = [];
      for (let H of X) {
        let B = { ...$, common: { ...$.common, issues: [] }, parent: null }, z = H._parseSync({ data: $.data, path: $.path, parent: B });
        if (z.status === "valid") return z;
        else if (z.status === "dirty" && !J) J = { result: z, ctx: B };
        if (B.common.issues.length) G.push(B.common.issues);
      }
      if (J) return $.common.issues.push(...J.ctx.common.issues), J.result;
      let W = G.map((H) => new x0(H));
      return E($, { code: A.invalid_union, unionErrors: W }), x;
    }
  }
  get options() {
    return this._def.options;
  }
};
X4.create = (Q, $) => {
  return new X4({ options: Q, typeName: j.ZodUnion, ...m($) });
};
var C1 = (Q) => {
  if (Q instanceof J4) return C1(Q.schema);
  else if (Q instanceof G1) return C1(Q.innerType());
  else if (Q instanceof G4) return [Q.value];
  else if (Q instanceof d1) return Q.options;
  else if (Q instanceof W4) return d.objectValues(Q.enum);
  else if (Q instanceof H4) return C1(Q._def.innerType);
  else if (Q instanceof Q4) return [void 0];
  else if (Q instanceof $4) return [null];
  else if (Q instanceof m0) return [void 0, ...C1(Q.unwrap())];
  else if (Q instanceof S1) return [null, ...C1(Q.unwrap())];
  else if (Q instanceof W8) return C1(Q.unwrap());
  else if (Q instanceof z4) return C1(Q.unwrap());
  else if (Q instanceof B4) return C1(Q._def.innerType);
  else return [];
};
var G8 = class _G8 extends p {
  _parse(Q) {
    let { ctx: $ } = this._processInputParams(Q);
    if ($.parsedType !== I.object) return E($, { code: A.invalid_type, expected: I.object, received: $.parsedType }), x;
    let X = this.discriminator, Y = $.data[X], J = this.optionsMap.get(Y);
    if (!J) return E($, { code: A.invalid_union_discriminator, options: Array.from(this.optionsMap.keys()), path: [X] }), x;
    if ($.common.async) return J._parseAsync({ data: $.data, path: $.path, parent: $ });
    else return J._parseSync({ data: $.data, path: $.path, parent: $ });
  }
  get discriminator() {
    return this._def.discriminator;
  }
  get options() {
    return this._def.options;
  }
  get optionsMap() {
    return this._def.optionsMap;
  }
  static create(Q, $, X) {
    let Y = /* @__PURE__ */ new Map();
    for (let J of $) {
      let G = C1(J.shape[Q]);
      if (!G.length) throw Error(`A discriminator value for key \`${Q}\` could not be extracted from all schema options`);
      for (let W of G) {
        if (Y.has(W)) throw Error(`Discriminator property ${String(Q)} has duplicate value ${String(W)}`);
        Y.set(W, J);
      }
    }
    return new _G8({ typeName: j.ZodDiscriminatedUnion, discriminator: Q, options: $, optionsMap: Y, ...m(X) });
  }
};
function B$(Q, $) {
  let X = w1(Q), Y = w1($);
  if (Q === $) return { valid: true, data: Q };
  else if (X === I.object && Y === I.object) {
    let J = d.objectKeys($), G = d.objectKeys(Q).filter((H) => J.indexOf(H) !== -1), W = { ...Q, ...$ };
    for (let H of G) {
      let B = B$(Q[H], $[H]);
      if (!B.valid) return { valid: false };
      W[H] = B.data;
    }
    return { valid: true, data: W };
  } else if (X === I.array && Y === I.array) {
    if (Q.length !== $.length) return { valid: false };
    let J = [];
    for (let G = 0; G < Q.length; G++) {
      let W = Q[G], H = $[G], B = B$(W, H);
      if (!B.valid) return { valid: false };
      J.push(B.data);
    }
    return { valid: true, data: J };
  } else if (X === I.date && Y === I.date && +Q === +$) return { valid: true, data: Q };
  else return { valid: false };
}
var Y4 = class extends p {
  _parse(Q) {
    let { status: $, ctx: X } = this._processInputParams(Q), Y = (J, G) => {
      if (Y8(J) || Y8(G)) return x;
      let W = B$(J.value, G.value);
      if (!W.valid) return E(X, { code: A.invalid_intersection_types }), x;
      if (J8(J) || J8(G)) $.dirty();
      return { status: $.value, value: W.data };
    };
    if (X.common.async) return Promise.all([this._def.left._parseAsync({ data: X.data, path: X.path, parent: X }), this._def.right._parseAsync({ data: X.data, path: X.path, parent: X })]).then(([J, G]) => Y(J, G));
    else return Y(this._def.left._parseSync({ data: X.data, path: X.path, parent: X }), this._def.right._parseSync({ data: X.data, path: X.path, parent: X }));
  }
};
Y4.create = (Q, $, X) => {
  return new Y4({ left: Q, right: $, typeName: j.ZodIntersection, ...m(X) });
};
var A1 = class _A1 extends p {
  _parse(Q) {
    let { status: $, ctx: X } = this._processInputParams(Q);
    if (X.parsedType !== I.array) return E(X, { code: A.invalid_type, expected: I.array, received: X.parsedType }), x;
    if (X.data.length < this._def.items.length) return E(X, { code: A.too_small, minimum: this._def.items.length, inclusive: true, exact: false, type: "array" }), x;
    if (!this._def.rest && X.data.length > this._def.items.length) E(X, { code: A.too_big, maximum: this._def.items.length, inclusive: true, exact: false, type: "array" }), $.dirty();
    let J = [...X.data].map((G, W) => {
      let H = this._def.items[W] || this._def.rest;
      if (!H) return null;
      return H._parse(new J1(X, G, X.path, W));
    }).filter((G) => !!G);
    if (X.common.async) return Promise.all(J).then((G) => {
      return A0.mergeArray($, G);
    });
    else return A0.mergeArray($, J);
  }
  get items() {
    return this._def.items;
  }
  rest(Q) {
    return new _A1({ ...this._def, rest: Q });
  }
};
A1.create = (Q, $) => {
  if (!Array.isArray(Q)) throw Error("You must pass an array of schemas to z.tuple([ ... ])");
  return new A1({ items: Q, typeName: j.ZodTuple, rest: null, ...m($) });
};
var d4 = class _d4 extends p {
  get keySchema() {
    return this._def.keyType;
  }
  get valueSchema() {
    return this._def.valueType;
  }
  _parse(Q) {
    let { status: $, ctx: X } = this._processInputParams(Q);
    if (X.parsedType !== I.object) return E(X, { code: A.invalid_type, expected: I.object, received: X.parsedType }), x;
    let Y = [], J = this._def.keyType, G = this._def.valueType;
    for (let W in X.data) Y.push({ key: J._parse(new J1(X, W, X.path, W)), value: G._parse(new J1(X, X.data[W], X.path, W)), alwaysSet: W in X.data });
    if (X.common.async) return A0.mergeObjectAsync($, Y);
    else return A0.mergeObjectSync($, Y);
  }
  get element() {
    return this._def.valueType;
  }
  static create(Q, $, X) {
    if ($ instanceof p) return new _d4({ keyType: Q, valueType: $, typeName: j.ZodRecord, ...m(X) });
    return new _d4({ keyType: X1.create(), valueType: Q, typeName: j.ZodRecord, ...m($) });
  }
};
var i4 = class extends p {
  get keySchema() {
    return this._def.keyType;
  }
  get valueSchema() {
    return this._def.valueType;
  }
  _parse(Q) {
    let { status: $, ctx: X } = this._processInputParams(Q);
    if (X.parsedType !== I.map) return E(X, { code: A.invalid_type, expected: I.map, received: X.parsedType }), x;
    let Y = this._def.keyType, J = this._def.valueType, G = [...X.data.entries()].map(([W, H], B) => {
      return { key: Y._parse(new J1(X, W, X.path, [B, "key"])), value: J._parse(new J1(X, H, X.path, [B, "value"])) };
    });
    if (X.common.async) {
      let W = /* @__PURE__ */ new Map();
      return Promise.resolve().then(async () => {
        for (let H of G) {
          let B = await H.key, z = await H.value;
          if (B.status === "aborted" || z.status === "aborted") return x;
          if (B.status === "dirty" || z.status === "dirty") $.dirty();
          W.set(B.value, z.value);
        }
        return { status: $.value, value: W };
      });
    } else {
      let W = /* @__PURE__ */ new Map();
      for (let H of G) {
        let { key: B, value: z } = H;
        if (B.status === "aborted" || z.status === "aborted") return x;
        if (B.status === "dirty" || z.status === "dirty") $.dirty();
        W.set(B.value, z.value);
      }
      return { status: $.value, value: W };
    }
  }
};
i4.create = (Q, $, X) => {
  return new i4({ valueType: $, keyType: Q, typeName: j.ZodMap, ...m(X) });
};
var D6 = class _D6 extends p {
  _parse(Q) {
    let { status: $, ctx: X } = this._processInputParams(Q);
    if (X.parsedType !== I.set) return E(X, { code: A.invalid_type, expected: I.set, received: X.parsedType }), x;
    let Y = this._def;
    if (Y.minSize !== null) {
      if (X.data.size < Y.minSize.value) E(X, { code: A.too_small, minimum: Y.minSize.value, type: "set", inclusive: true, exact: false, message: Y.minSize.message }), $.dirty();
    }
    if (Y.maxSize !== null) {
      if (X.data.size > Y.maxSize.value) E(X, { code: A.too_big, maximum: Y.maxSize.value, type: "set", inclusive: true, exact: false, message: Y.maxSize.message }), $.dirty();
    }
    let J = this._def.valueType;
    function G(H) {
      let B = /* @__PURE__ */ new Set();
      for (let z of H) {
        if (z.status === "aborted") return x;
        if (z.status === "dirty") $.dirty();
        B.add(z.value);
      }
      return { status: $.value, value: B };
    }
    let W = [...X.data.values()].map((H, B) => J._parse(new J1(X, H, X.path, B)));
    if (X.common.async) return Promise.all(W).then((H) => G(H));
    else return G(W);
  }
  min(Q, $) {
    return new _D6({ ...this._def, minSize: { value: Q, message: C.toString($) } });
  }
  max(Q, $) {
    return new _D6({ ...this._def, maxSize: { value: Q, message: C.toString($) } });
  }
  size(Q, $) {
    return this.min(Q, $).max(Q, $);
  }
  nonempty(Q) {
    return this.min(1, Q);
  }
};
D6.create = (Q, $) => {
  return new D6({ valueType: Q, minSize: null, maxSize: null, typeName: j.ZodSet, ...m($) });
};
var s6 = class _s6 extends p {
  constructor() {
    super(...arguments);
    this.validate = this.implement;
  }
  _parse(Q) {
    let { ctx: $ } = this._processInputParams(Q);
    if ($.parsedType !== I.function) return E($, { code: A.invalid_type, expected: I.function, received: $.parsedType }), x;
    function X(W, H) {
      return l4({ data: W, path: $.path, errorMaps: [$.common.contextualErrorMap, $.schemaErrorMap, r6(), Z1].filter((B) => !!B), issueData: { code: A.invalid_arguments, argumentsError: H } });
    }
    function Y(W, H) {
      return l4({ data: W, path: $.path, errorMaps: [$.common.contextualErrorMap, $.schemaErrorMap, r6(), Z1].filter((B) => !!B), issueData: { code: A.invalid_return_type, returnTypeError: H } });
    }
    let J = { errorMap: $.common.contextualErrorMap }, G = $.data;
    if (this._def.returns instanceof w6) {
      let W = this;
      return b0(async function(...H) {
        let B = new x0([]), z = await W._def.args.parseAsync(H, J).catch((q) => {
          throw B.addIssue(X(H, q)), B;
        }), K = await Reflect.apply(G, this, z);
        return await W._def.returns._def.type.parseAsync(K, J).catch((q) => {
          throw B.addIssue(Y(K, q)), B;
        });
      });
    } else {
      let W = this;
      return b0(function(...H) {
        let B = W._def.args.safeParse(H, J);
        if (!B.success) throw new x0([X(H, B.error)]);
        let z = Reflect.apply(G, this, B.data), K = W._def.returns.safeParse(z, J);
        if (!K.success) throw new x0([Y(z, K.error)]);
        return K.data;
      });
    }
  }
  parameters() {
    return this._def.args;
  }
  returnType() {
    return this._def.returns;
  }
  args(...Q) {
    return new _s6({ ...this._def, args: A1.create(Q).rest(l1.create()) });
  }
  returns(Q) {
    return new _s6({ ...this._def, returns: Q });
  }
  implement(Q) {
    return this.parse(Q);
  }
  strictImplement(Q) {
    return this.parse(Q);
  }
  static create(Q, $, X) {
    return new _s6({ args: Q ? Q : A1.create([]).rest(l1.create()), returns: $ || l1.create(), typeName: j.ZodFunction, ...m(X) });
  }
};
var J4 = class extends p {
  get schema() {
    return this._def.getter();
  }
  _parse(Q) {
    let { ctx: $ } = this._processInputParams(Q);
    return this._def.getter()._parse({ data: $.data, path: $.path, parent: $ });
  }
};
J4.create = (Q, $) => {
  return new J4({ getter: Q, typeName: j.ZodLazy, ...m($) });
};
var G4 = class extends p {
  _parse(Q) {
    if (Q.data !== this._def.value) {
      let $ = this._getOrReturnCtx(Q);
      return E($, { received: $.data, code: A.invalid_literal, expected: this._def.value }), x;
    }
    return { status: "valid", value: Q.data };
  }
  get value() {
    return this._def.value;
  }
};
G4.create = (Q, $) => {
  return new G4({ value: Q, typeName: j.ZodLiteral, ...m($) });
};
function RJ(Q, $) {
  return new d1({ values: Q, typeName: j.ZodEnum, ...m($) });
}
var d1 = class _d1 extends p {
  _parse(Q) {
    if (typeof Q.data !== "string") {
      let $ = this._getOrReturnCtx(Q), X = this._def.values;
      return E($, { expected: d.joinValues(X), received: $.parsedType, code: A.invalid_type }), x;
    }
    if (!this._cache) this._cache = new Set(this._def.values);
    if (!this._cache.has(Q.data)) {
      let $ = this._getOrReturnCtx(Q), X = this._def.values;
      return E($, { received: $.data, code: A.invalid_enum_value, options: X }), x;
    }
    return b0(Q.data);
  }
  get options() {
    return this._def.values;
  }
  get enum() {
    let Q = {};
    for (let $ of this._def.values) Q[$] = $;
    return Q;
  }
  get Values() {
    let Q = {};
    for (let $ of this._def.values) Q[$] = $;
    return Q;
  }
  get Enum() {
    let Q = {};
    for (let $ of this._def.values) Q[$] = $;
    return Q;
  }
  extract(Q, $ = this._def) {
    return _d1.create(Q, { ...this._def, ...$ });
  }
  exclude(Q, $ = this._def) {
    return _d1.create(this.options.filter((X) => !Q.includes(X)), { ...this._def, ...$ });
  }
};
d1.create = RJ;
var W4 = class extends p {
  _parse(Q) {
    let $ = d.getValidEnumValues(this._def.values), X = this._getOrReturnCtx(Q);
    if (X.parsedType !== I.string && X.parsedType !== I.number) {
      let Y = d.objectValues($);
      return E(X, { expected: d.joinValues(Y), received: X.parsedType, code: A.invalid_type }), x;
    }
    if (!this._cache) this._cache = new Set(d.getValidEnumValues(this._def.values));
    if (!this._cache.has(Q.data)) {
      let Y = d.objectValues($);
      return E(X, { received: X.data, code: A.invalid_enum_value, options: Y }), x;
    }
    return b0(Q.data);
  }
  get enum() {
    return this._def.values;
  }
};
W4.create = (Q, $) => {
  return new W4({ values: Q, typeName: j.ZodNativeEnum, ...m($) });
};
var w6 = class extends p {
  unwrap() {
    return this._def.type;
  }
  _parse(Q) {
    let { ctx: $ } = this._processInputParams(Q);
    if ($.parsedType !== I.promise && $.common.async === false) return E($, { code: A.invalid_type, expected: I.promise, received: $.parsedType }), x;
    let X = $.parsedType === I.promise ? $.data : Promise.resolve($.data);
    return b0(X.then((Y) => {
      return this._def.type.parseAsync(Y, { path: $.path, errorMap: $.common.contextualErrorMap });
    }));
  }
};
w6.create = (Q, $) => {
  return new w6({ type: Q, typeName: j.ZodPromise, ...m($) });
};
var G1 = class extends p {
  innerType() {
    return this._def.schema;
  }
  sourceType() {
    return this._def.schema._def.typeName === j.ZodEffects ? this._def.schema.sourceType() : this._def.schema;
  }
  _parse(Q) {
    let { status: $, ctx: X } = this._processInputParams(Q), Y = this._def.effect || null, J = { addIssue: (G) => {
      if (E(X, G), G.fatal) $.abort();
      else $.dirty();
    }, get path() {
      return X.path;
    } };
    if (J.addIssue = J.addIssue.bind(J), Y.type === "preprocess") {
      let G = Y.transform(X.data, J);
      if (X.common.async) return Promise.resolve(G).then(async (W) => {
        if ($.value === "aborted") return x;
        let H = await this._def.schema._parseAsync({ data: W, path: X.path, parent: X });
        if (H.status === "aborted") return x;
        if (H.status === "dirty") return F6(H.value);
        if ($.value === "dirty") return F6(H.value);
        return H;
      });
      else {
        if ($.value === "aborted") return x;
        let W = this._def.schema._parseSync({ data: G, path: X.path, parent: X });
        if (W.status === "aborted") return x;
        if (W.status === "dirty") return F6(W.value);
        if ($.value === "dirty") return F6(W.value);
        return W;
      }
    }
    if (Y.type === "refinement") {
      let G = (W) => {
        let H = Y.refinement(W, J);
        if (X.common.async) return Promise.resolve(H);
        if (H instanceof Promise) throw Error("Async refinement encountered during synchronous parse operation. Use .parseAsync instead.");
        return W;
      };
      if (X.common.async === false) {
        let W = this._def.schema._parseSync({ data: X.data, path: X.path, parent: X });
        if (W.status === "aborted") return x;
        if (W.status === "dirty") $.dirty();
        return G(W.value), { status: $.value, value: W.value };
      } else return this._def.schema._parseAsync({ data: X.data, path: X.path, parent: X }).then((W) => {
        if (W.status === "aborted") return x;
        if (W.status === "dirty") $.dirty();
        return G(W.value).then(() => {
          return { status: $.value, value: W.value };
        });
      });
    }
    if (Y.type === "transform") if (X.common.async === false) {
      let G = this._def.schema._parseSync({ data: X.data, path: X.path, parent: X });
      if (!m1(G)) return x;
      let W = Y.transform(G.value, J);
      if (W instanceof Promise) throw Error("Asynchronous transform encountered during synchronous parse operation. Use .parseAsync instead.");
      return { status: $.value, value: W };
    } else return this._def.schema._parseAsync({ data: X.data, path: X.path, parent: X }).then((G) => {
      if (!m1(G)) return x;
      return Promise.resolve(Y.transform(G.value, J)).then((W) => ({ status: $.value, value: W }));
    });
    d.assertNever(Y);
  }
};
G1.create = (Q, $, X) => {
  return new G1({ schema: Q, typeName: j.ZodEffects, effect: $, ...m(X) });
};
G1.createWithPreprocess = (Q, $, X) => {
  return new G1({ schema: $, effect: { type: "preprocess", transform: Q }, typeName: j.ZodEffects, ...m(X) });
};
var m0 = class extends p {
  _parse(Q) {
    if (this._getType(Q) === I.undefined) return b0(void 0);
    return this._def.innerType._parse(Q);
  }
  unwrap() {
    return this._def.innerType;
  }
};
m0.create = (Q, $) => {
  return new m0({ innerType: Q, typeName: j.ZodOptional, ...m($) });
};
var S1 = class extends p {
  _parse(Q) {
    if (this._getType(Q) === I.null) return b0(null);
    return this._def.innerType._parse(Q);
  }
  unwrap() {
    return this._def.innerType;
  }
};
S1.create = (Q, $) => {
  return new S1({ innerType: Q, typeName: j.ZodNullable, ...m($) });
};
var H4 = class extends p {
  _parse(Q) {
    let { ctx: $ } = this._processInputParams(Q), X = $.data;
    if ($.parsedType === I.undefined) X = this._def.defaultValue();
    return this._def.innerType._parse({ data: X, path: $.path, parent: $ });
  }
  removeDefault() {
    return this._def.innerType;
  }
};
H4.create = (Q, $) => {
  return new H4({ innerType: Q, typeName: j.ZodDefault, defaultValue: typeof $.default === "function" ? $.default : () => $.default, ...m($) });
};
var B4 = class extends p {
  _parse(Q) {
    let { ctx: $ } = this._processInputParams(Q), X = { ...$, common: { ...$.common, issues: [] } }, Y = this._def.innerType._parse({ data: X.data, path: X.path, parent: { ...X } });
    if (t6(Y)) return Y.then((J) => {
      return { status: "valid", value: J.status === "valid" ? J.value : this._def.catchValue({ get error() {
        return new x0(X.common.issues);
      }, input: X.data }) };
    });
    else return { status: "valid", value: Y.status === "valid" ? Y.value : this._def.catchValue({ get error() {
      return new x0(X.common.issues);
    }, input: X.data }) };
  }
  removeCatch() {
    return this._def.innerType;
  }
};
B4.create = (Q, $) => {
  return new B4({ innerType: Q, typeName: j.ZodCatch, catchValue: typeof $.catch === "function" ? $.catch : () => $.catch, ...m($) });
};
var n4 = class extends p {
  _parse(Q) {
    if (this._getType(Q) !== I.nan) {
      let X = this._getOrReturnCtx(Q);
      return E(X, { code: A.invalid_type, expected: I.nan, received: X.parsedType }), x;
    }
    return { status: "valid", value: Q.data };
  }
};
n4.create = (Q) => {
  return new n4({ typeName: j.ZodNaN, ...m(Q) });
};
var tL = /* @__PURE__ */ Symbol("zod_brand");
var W8 = class extends p {
  _parse(Q) {
    let { ctx: $ } = this._processInputParams(Q), X = $.data;
    return this._def.type._parse({ data: X, path: $.path, parent: $ });
  }
  unwrap() {
    return this._def.type;
  }
};
var o4 = class _o4 extends p {
  _parse(Q) {
    let { status: $, ctx: X } = this._processInputParams(Q);
    if (X.common.async) return (async () => {
      let J = await this._def.in._parseAsync({ data: X.data, path: X.path, parent: X });
      if (J.status === "aborted") return x;
      if (J.status === "dirty") return $.dirty(), F6(J.value);
      else return this._def.out._parseAsync({ data: J.value, path: X.path, parent: X });
    })();
    else {
      let Y = this._def.in._parseSync({ data: X.data, path: X.path, parent: X });
      if (Y.status === "aborted") return x;
      if (Y.status === "dirty") return $.dirty(), { status: "dirty", value: Y.value };
      else return this._def.out._parseSync({ data: Y.value, path: X.path, parent: X });
    }
  }
  static create(Q, $) {
    return new _o4({ in: Q, out: $, typeName: j.ZodPipeline });
  }
};
var z4 = class extends p {
  _parse(Q) {
    let $ = this._def.innerType._parse(Q), X = (Y) => {
      if (m1(Y)) Y.value = Object.freeze(Y.value);
      return Y;
    };
    return t6($) ? $.then((Y) => X(Y)) : X($);
  }
  unwrap() {
    return this._def.innerType;
  }
};
z4.create = (Q, $) => {
  return new z4({ innerType: Q, typeName: j.ZodReadonly, ...m($) });
};
function wJ(Q, $) {
  let X = typeof Q === "function" ? Q($) : typeof Q === "string" ? { message: Q } : Q;
  return typeof X === "string" ? { message: X } : X;
}
function IJ(Q, $ = {}, X) {
  if (Q) return O6.create().superRefine((Y, J) => {
    let G = Q(Y);
    if (G instanceof Promise) return G.then((W) => {
      if (!W) {
        let H = wJ($, Y), B = H.fatal ?? X ?? true;
        J.addIssue({ code: "custom", ...H, fatal: B });
      }
    });
    if (!G) {
      let W = wJ($, Y), H = W.fatal ?? X ?? true;
      J.addIssue({ code: "custom", ...W, fatal: H });
    }
    return;
  });
  return O6.create();
}
var aL = { object: U0.lazycreate };
var j;
(function(Q) {
  Q.ZodString = "ZodString", Q.ZodNumber = "ZodNumber", Q.ZodNaN = "ZodNaN", Q.ZodBigInt = "ZodBigInt", Q.ZodBoolean = "ZodBoolean", Q.ZodDate = "ZodDate", Q.ZodSymbol = "ZodSymbol", Q.ZodUndefined = "ZodUndefined", Q.ZodNull = "ZodNull", Q.ZodAny = "ZodAny", Q.ZodUnknown = "ZodUnknown", Q.ZodNever = "ZodNever", Q.ZodVoid = "ZodVoid", Q.ZodArray = "ZodArray", Q.ZodObject = "ZodObject", Q.ZodUnion = "ZodUnion", Q.ZodDiscriminatedUnion = "ZodDiscriminatedUnion", Q.ZodIntersection = "ZodIntersection", Q.ZodTuple = "ZodTuple", Q.ZodRecord = "ZodRecord", Q.ZodMap = "ZodMap", Q.ZodSet = "ZodSet", Q.ZodFunction = "ZodFunction", Q.ZodLazy = "ZodLazy", Q.ZodLiteral = "ZodLiteral", Q.ZodEnum = "ZodEnum", Q.ZodEffects = "ZodEffects", Q.ZodNativeEnum = "ZodNativeEnum", Q.ZodOptional = "ZodOptional", Q.ZodNullable = "ZodNullable", Q.ZodDefault = "ZodDefault", Q.ZodCatch = "ZodCatch", Q.ZodPromise = "ZodPromise", Q.ZodBranded = "ZodBranded", Q.ZodPipeline = "ZodPipeline", Q.ZodReadonly = "ZodReadonly";
})(j || (j = {}));
var sL = (Q, $ = { message: `Input not instance of ${Q.name}` }) => IJ((X) => X instanceof Q, $);
var EJ = X1.create;
var PJ = c1.create;
var eL = n4.create;
var QF = p1.create;
var bJ = e6.create;
var $F = N6.create;
var XF = c4.create;
var YF = Q4.create;
var JF = $4.create;
var GF = O6.create;
var WF = l1.create;
var HF = M1.create;
var BF = p4.create;
var zF = Y1.create;
var z$ = U0.create;
var KF = U0.strictCreate;
var VF = X4.create;
var qF = G8.create;
var UF = Y4.create;
var LF = A1.create;
var FF = d4.create;
var NF = i4.create;
var OF = D6.create;
var DF = s6.create;
var wF = J4.create;
var MF = G4.create;
var AF = d1.create;
var jF = W4.create;
var RF = w6.create;
var IF = G1.create;
var EF = m0.create;
var PF = S1.create;
var bF = G1.createWithPreprocess;
var ZF = o4.create;
var CF = () => EJ().optional();
var SF = () => PJ().optional();
var _F = () => bJ().optional();
var kF = { string: (Q) => X1.create({ ...Q, coerce: true }), number: (Q) => c1.create({ ...Q, coerce: true }), boolean: (Q) => e6.create({ ...Q, coerce: true }), bigint: (Q) => p1.create({ ...Q, coerce: true }), date: (Q) => N6.create({ ...Q, coerce: true }) };
var vF = x;
var TF = Object.freeze({ status: "aborted" });
function O(Q, $, X) {
  function Y(H, B) {
    var z;
    Object.defineProperty(H, "_zod", { value: H._zod ?? {}, enumerable: false }), (z = H._zod).traits ?? (z.traits = /* @__PURE__ */ new Set()), H._zod.traits.add(Q), $(H, B);
    for (let K in W.prototype) if (!(K in H)) Object.defineProperty(H, K, { value: W.prototype[K].bind(H) });
    H._zod.constr = W, H._zod.def = B;
  }
  let J = X?.Parent ?? Object;
  class G extends J {
  }
  Object.defineProperty(G, "name", { value: Q });
  function W(H) {
    var B;
    let z = X?.Parent ? new G() : this;
    Y(z, H), (B = z._zod).deferred ?? (B.deferred = []);
    for (let K of z._zod.deferred) K();
    return z;
  }
  return Object.defineProperty(W, "init", { value: Y }), Object.defineProperty(W, Symbol.hasInstance, { value: (H) => {
    if (X?.Parent && H instanceof X.Parent) return true;
    return H?._zod?.traits?.has(Q);
  } }), Object.defineProperty(W, "name", { value: Q }), W;
}
var n1 = class extends Error {
  constructor() {
    super("Encountered Promise during synchronous parse. Use .parseAsync() instead.");
  }
};
var H8 = {};
function l0(Q) {
  if (Q) Object.assign(H8, Q);
  return H8;
}
var i = {};
uQ(i, { unwrapMessage: () => r4, stringifyPrimitive: () => K8, required: () => sF, randomString: () => cF, propertyKeyTypes: () => F$, promiseAllObject: () => lF, primitiveTypes: () => ZJ, prefixIssues: () => j1, pick: () => nF, partial: () => aF, optionalKeys: () => N$, omit: () => oF, numKeys: () => pF, nullish: () => s4, normalizeParams: () => y, merge: () => tF, jsonStringifyReplacer: () => V$, joinValues: () => B8, issue: () => D$, isPlainObject: () => V4, isObject: () => K4, getSizableOrigin: () => SJ, getParsedType: () => dF, getLengthableOrigin: () => Q9, getEnumValues: () => t4, getElementAtPath: () => mF, floatSafeRemainder: () => q$, finalizeIssue: () => W1, extend: () => rF, escapeRegex: () => o1, esc: () => M6, defineLazy: () => Y0, createTransparentProxy: () => iF, clone: () => c0, cleanRegex: () => e4, cleanEnum: () => eF, captureStackTrace: () => z8, cached: () => a4, assignProp: () => U$, assertNotEqual: () => gF, assertNever: () => fF, assertIs: () => hF, assertEqual: () => yF, assert: () => uF, allowsEval: () => L$, aborted: () => A6, NUMBER_FORMAT_RANGES: () => O$, Class: () => _J, BIGINT_FORMAT_RANGES: () => CJ });
function yF(Q) {
  return Q;
}
function gF(Q) {
  return Q;
}
function hF(Q) {
}
function fF(Q) {
  throw Error();
}
function uF(Q) {
}
function t4(Q) {
  let $ = Object.values(Q).filter((Y) => typeof Y === "number");
  return Object.entries(Q).filter(([Y, J]) => $.indexOf(+Y) === -1).map(([Y, J]) => J);
}
function B8(Q, $ = "|") {
  return Q.map((X) => K8(X)).join($);
}
function V$(Q, $) {
  if (typeof $ === "bigint") return $.toString();
  return $;
}
function a4(Q) {
  return { get value() {
    {
      let X = Q();
      return Object.defineProperty(this, "value", { value: X }), X;
    }
    throw Error("cached value already set");
  } };
}
function s4(Q) {
  return Q === null || Q === void 0;
}
function e4(Q) {
  let $ = Q.startsWith("^") ? 1 : 0, X = Q.endsWith("$") ? Q.length - 1 : Q.length;
  return Q.slice($, X);
}
function q$(Q, $) {
  let X = (Q.toString().split(".")[1] || "").length, Y = ($.toString().split(".")[1] || "").length, J = X > Y ? X : Y, G = Number.parseInt(Q.toFixed(J).replace(".", "")), W = Number.parseInt($.toFixed(J).replace(".", ""));
  return G % W / 10 ** J;
}
function Y0(Q, $, X) {
  Object.defineProperty(Q, $, { get() {
    {
      let J = X();
      return Q[$] = J, J;
    }
    throw Error("cached value already set");
  }, set(J) {
    Object.defineProperty(Q, $, { value: J });
  }, configurable: true });
}
function U$(Q, $, X) {
  Object.defineProperty(Q, $, { value: X, writable: true, enumerable: true, configurable: true });
}
function mF(Q, $) {
  if (!$) return Q;
  return $.reduce((X, Y) => X?.[Y], Q);
}
function lF(Q) {
  let $ = Object.keys(Q), X = $.map((Y) => Q[Y]);
  return Promise.all(X).then((Y) => {
    let J = {};
    for (let G = 0; G < $.length; G++) J[$[G]] = Y[G];
    return J;
  });
}
function cF(Q = 10) {
  let X = "";
  for (let Y = 0; Y < Q; Y++) X += "abcdefghijklmnopqrstuvwxyz"[Math.floor(Math.random() * 26)];
  return X;
}
function M6(Q) {
  return JSON.stringify(Q);
}
var z8 = Error.captureStackTrace ? Error.captureStackTrace : (...Q) => {
};
function K4(Q) {
  return typeof Q === "object" && Q !== null && !Array.isArray(Q);
}
var L$ = a4(() => {
  if (typeof navigator < "u" && navigator?.userAgent?.includes("Cloudflare")) return false;
  try {
    return new Function(""), true;
  } catch (Q) {
    return false;
  }
});
function V4(Q) {
  if (K4(Q) === false) return false;
  let $ = Q.constructor;
  if ($ === void 0) return true;
  let X = $.prototype;
  if (K4(X) === false) return false;
  if (Object.prototype.hasOwnProperty.call(X, "isPrototypeOf") === false) return false;
  return true;
}
function pF(Q) {
  let $ = 0;
  for (let X in Q) if (Object.prototype.hasOwnProperty.call(Q, X)) $++;
  return $;
}
var dF = (Q) => {
  let $ = typeof Q;
  switch ($) {
    case "undefined":
      return "undefined";
    case "string":
      return "string";
    case "number":
      return Number.isNaN(Q) ? "nan" : "number";
    case "boolean":
      return "boolean";
    case "function":
      return "function";
    case "bigint":
      return "bigint";
    case "symbol":
      return "symbol";
    case "object":
      if (Array.isArray(Q)) return "array";
      if (Q === null) return "null";
      if (Q.then && typeof Q.then === "function" && Q.catch && typeof Q.catch === "function") return "promise";
      if (typeof Map < "u" && Q instanceof Map) return "map";
      if (typeof Set < "u" && Q instanceof Set) return "set";
      if (typeof Date < "u" && Q instanceof Date) return "date";
      if (typeof File < "u" && Q instanceof File) return "file";
      return "object";
    default:
      throw Error(`Unknown data type: ${$}`);
  }
};
var F$ = /* @__PURE__ */ new Set(["string", "number", "symbol"]);
var ZJ = /* @__PURE__ */ new Set(["string", "number", "bigint", "boolean", "symbol", "undefined"]);
function o1(Q) {
  return Q.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}
function c0(Q, $, X) {
  let Y = new Q._zod.constr($ ?? Q._zod.def);
  if (!$ || X?.parent) Y._zod.parent = Q;
  return Y;
}
function y(Q) {
  let $ = Q;
  if (!$) return {};
  if (typeof $ === "string") return { error: () => $ };
  if ($?.message !== void 0) {
    if ($?.error !== void 0) throw Error("Cannot specify both `message` and `error` params");
    $.error = $.message;
  }
  if (delete $.message, typeof $.error === "string") return { ...$, error: () => $.error };
  return $;
}
function iF(Q) {
  let $;
  return new Proxy({}, { get(X, Y, J) {
    return $ ?? ($ = Q()), Reflect.get($, Y, J);
  }, set(X, Y, J, G) {
    return $ ?? ($ = Q()), Reflect.set($, Y, J, G);
  }, has(X, Y) {
    return $ ?? ($ = Q()), Reflect.has($, Y);
  }, deleteProperty(X, Y) {
    return $ ?? ($ = Q()), Reflect.deleteProperty($, Y);
  }, ownKeys(X) {
    return $ ?? ($ = Q()), Reflect.ownKeys($);
  }, getOwnPropertyDescriptor(X, Y) {
    return $ ?? ($ = Q()), Reflect.getOwnPropertyDescriptor($, Y);
  }, defineProperty(X, Y, J) {
    return $ ?? ($ = Q()), Reflect.defineProperty($, Y, J);
  } });
}
function K8(Q) {
  if (typeof Q === "bigint") return Q.toString() + "n";
  if (typeof Q === "string") return `"${Q}"`;
  return `${Q}`;
}
function N$(Q) {
  return Object.keys(Q).filter(($) => {
    return Q[$]._zod.optin === "optional" && Q[$]._zod.optout === "optional";
  });
}
var O$ = { safeint: [Number.MIN_SAFE_INTEGER, Number.MAX_SAFE_INTEGER], int32: [-2147483648, 2147483647], uint32: [0, 4294967295], float32: [-34028234663852886e22, 34028234663852886e22], float64: [-Number.MAX_VALUE, Number.MAX_VALUE] };
var CJ = { int64: [BigInt("-9223372036854775808"), BigInt("9223372036854775807")], uint64: [BigInt(0), BigInt("18446744073709551615")] };
function nF(Q, $) {
  let X = {}, Y = Q._zod.def;
  for (let J in $) {
    if (!(J in Y.shape)) throw Error(`Unrecognized key: "${J}"`);
    if (!$[J]) continue;
    X[J] = Y.shape[J];
  }
  return c0(Q, { ...Q._zod.def, shape: X, checks: [] });
}
function oF(Q, $) {
  let X = { ...Q._zod.def.shape }, Y = Q._zod.def;
  for (let J in $) {
    if (!(J in Y.shape)) throw Error(`Unrecognized key: "${J}"`);
    if (!$[J]) continue;
    delete X[J];
  }
  return c0(Q, { ...Q._zod.def, shape: X, checks: [] });
}
function rF(Q, $) {
  if (!V4($)) throw Error("Invalid input to extend: expected a plain object");
  let X = { ...Q._zod.def, get shape() {
    let Y = { ...Q._zod.def.shape, ...$ };
    return U$(this, "shape", Y), Y;
  }, checks: [] };
  return c0(Q, X);
}
function tF(Q, $) {
  return c0(Q, { ...Q._zod.def, get shape() {
    let X = { ...Q._zod.def.shape, ...$._zod.def.shape };
    return U$(this, "shape", X), X;
  }, catchall: $._zod.def.catchall, checks: [] });
}
function aF(Q, $, X) {
  let Y = $._zod.def.shape, J = { ...Y };
  if (X) for (let G in X) {
    if (!(G in Y)) throw Error(`Unrecognized key: "${G}"`);
    if (!X[G]) continue;
    J[G] = Q ? new Q({ type: "optional", innerType: Y[G] }) : Y[G];
  }
  else for (let G in Y) J[G] = Q ? new Q({ type: "optional", innerType: Y[G] }) : Y[G];
  return c0($, { ...$._zod.def, shape: J, checks: [] });
}
function sF(Q, $, X) {
  let Y = $._zod.def.shape, J = { ...Y };
  if (X) for (let G in X) {
    if (!(G in J)) throw Error(`Unrecognized key: "${G}"`);
    if (!X[G]) continue;
    J[G] = new Q({ type: "nonoptional", innerType: Y[G] });
  }
  else for (let G in Y) J[G] = new Q({ type: "nonoptional", innerType: Y[G] });
  return c0($, { ...$._zod.def, shape: J, checks: [] });
}
function A6(Q, $ = 0) {
  for (let X = $; X < Q.issues.length; X++) if (Q.issues[X]?.continue !== true) return true;
  return false;
}
function j1(Q, $) {
  return $.map((X) => {
    var Y;
    return (Y = X).path ?? (Y.path = []), X.path.unshift(Q), X;
  });
}
function r4(Q) {
  return typeof Q === "string" ? Q : Q?.message;
}
function W1(Q, $, X) {
  let Y = { ...Q, path: Q.path ?? [] };
  if (!Q.message) {
    let J = r4(Q.inst?._zod.def?.error?.(Q)) ?? r4($?.error?.(Q)) ?? r4(X.customError?.(Q)) ?? r4(X.localeError?.(Q)) ?? "Invalid input";
    Y.message = J;
  }
  if (delete Y.inst, delete Y.continue, !$?.reportInput) delete Y.input;
  return Y;
}
function SJ(Q) {
  if (Q instanceof Set) return "set";
  if (Q instanceof Map) return "map";
  if (Q instanceof File) return "file";
  return "unknown";
}
function Q9(Q) {
  if (Array.isArray(Q)) return "array";
  if (typeof Q === "string") return "string";
  return "unknown";
}
function D$(...Q) {
  let [$, X, Y] = Q;
  if (typeof $ === "string") return { message: $, code: "custom", input: X, inst: Y };
  return { ...$ };
}
function eF(Q) {
  return Object.entries(Q).filter(([$, X]) => {
    return Number.isNaN(Number.parseInt($, 10));
  }).map(($) => $[1]);
}
var _J = class {
  constructor(...Q) {
  }
};
var kJ = (Q, $) => {
  Q.name = "$ZodError", Object.defineProperty(Q, "_zod", { value: Q._zod, enumerable: false }), Object.defineProperty(Q, "issues", { value: $, enumerable: false }), Object.defineProperty(Q, "message", { get() {
    return JSON.stringify($, V$, 2);
  }, enumerable: true });
};
var V8 = O("$ZodError", kJ);
var $9 = O("$ZodError", kJ, { Parent: Error });
function w$(Q, $ = (X) => X.message) {
  let X = {}, Y = [];
  for (let J of Q.issues) if (J.path.length > 0) X[J.path[0]] = X[J.path[0]] || [], X[J.path[0]].push($(J));
  else Y.push($(J));
  return { formErrors: Y, fieldErrors: X };
}
function M$(Q, $) {
  let X = $ || function(G) {
    return G.message;
  }, Y = { _errors: [] }, J = (G) => {
    for (let W of G.issues) if (W.code === "invalid_union" && W.errors.length) W.errors.map((H) => J({ issues: H }));
    else if (W.code === "invalid_key") J({ issues: W.issues });
    else if (W.code === "invalid_element") J({ issues: W.issues });
    else if (W.path.length === 0) Y._errors.push(X(W));
    else {
      let H = Y, B = 0;
      while (B < W.path.length) {
        let z = W.path[B];
        if (B !== W.path.length - 1) H[z] = H[z] || { _errors: [] };
        else H[z] = H[z] || { _errors: [] }, H[z]._errors.push(X(W));
        H = H[z], B++;
      }
    }
  };
  return J(Q), Y;
}
var A$ = (Q) => ($, X, Y, J) => {
  let G = Y ? Object.assign(Y, { async: false }) : { async: false }, W = $._zod.run({ value: X, issues: [] }, G);
  if (W instanceof Promise) throw new n1();
  if (W.issues.length) {
    let H = new (J?.Err ?? Q)(W.issues.map((B) => W1(B, G, l0())));
    throw z8(H, J?.callee), H;
  }
  return W.value;
};
var j$ = A$($9);
var R$ = (Q) => async ($, X, Y, J) => {
  let G = Y ? Object.assign(Y, { async: true }) : { async: true }, W = $._zod.run({ value: X, issues: [] }, G);
  if (W instanceof Promise) W = await W;
  if (W.issues.length) {
    let H = new (J?.Err ?? Q)(W.issues.map((B) => W1(B, G, l0())));
    throw z8(H, J?.callee), H;
  }
  return W.value;
};
var I$ = R$($9);
var E$ = (Q) => ($, X, Y) => {
  let J = Y ? { ...Y, async: false } : { async: false }, G = $._zod.run({ value: X, issues: [] }, J);
  if (G instanceof Promise) throw new n1();
  return G.issues.length ? { success: false, error: new (Q ?? V8)(G.issues.map((W) => W1(W, J, l0()))) } : { success: true, data: G.value };
};
var j6 = E$($9);
var P$ = (Q) => async ($, X, Y) => {
  let J = Y ? Object.assign(Y, { async: true }) : { async: true }, G = $._zod.run({ value: X, issues: [] }, J);
  if (G instanceof Promise) G = await G;
  return G.issues.length ? { success: false, error: new Q(G.issues.map((W) => W1(W, J, l0()))) } : { success: true, data: G.value };
};
var R6 = P$($9);
var vJ = /^[cC][^\s-]{8,}$/;
var TJ = /^[0-9a-z]+$/;
var xJ = /^[0-9A-HJKMNP-TV-Za-hjkmnp-tv-z]{26}$/;
var yJ = /^[0-9a-vA-V]{20}$/;
var gJ = /^[A-Za-z0-9]{27}$/;
var hJ = /^[a-zA-Z0-9_-]{21}$/;
var fJ = /^P(?:(\d+W)|(?!.*W)(?=\d|T\d)(\d+Y)?(\d+M)?(\d+D)?(T(?=\d)(\d+H)?(\d+M)?(\d+([.,]\d+)?S)?)?)$/;
var uJ = /^([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$/;
var b$ = (Q) => {
  if (!Q) return /^([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-8][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}|00000000-0000-0000-0000-000000000000)$/;
  return new RegExp(`^([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-${Q}[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12})$`);
};
var mJ = /^(?!\.)(?!.*\.\.)([A-Za-z0-9_'+\-\.]*)[A-Za-z0-9_+-]@([A-Za-z0-9][A-Za-z0-9\-]*\.)+[A-Za-z]{2,}$/;
function lJ() {
  return new RegExp("^(\\p{Extended_Pictographic}|\\p{Emoji_Component})+$", "u");
}
var cJ = /^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])$/;
var pJ = /^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|::|([0-9a-fA-F]{1,4})?::([0-9a-fA-F]{1,4}:?){0,6})$/;
var dJ = /^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\/([0-9]|[1-2][0-9]|3[0-2])$/;
var iJ = /^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|::|([0-9a-fA-F]{1,4})?::([0-9a-fA-F]{1,4}:?){0,6})\/(12[0-8]|1[01][0-9]|[1-9]?[0-9])$/;
var nJ = /^$|^(?:[0-9a-zA-Z+/]{4})*(?:(?:[0-9a-zA-Z+/]{2}==)|(?:[0-9a-zA-Z+/]{3}=))?$/;
var Z$ = /^[A-Za-z0-9_-]*$/;
var oJ = /^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+$/;
var rJ = /^\+(?:[0-9]){6,14}[0-9]$/;
var tJ = "(?:(?:\\d\\d[2468][048]|\\d\\d[13579][26]|\\d\\d0[48]|[02468][048]00|[13579][26]00)-02-29|\\d{4}-(?:(?:0[13578]|1[02])-(?:0[1-9]|[12]\\d|3[01])|(?:0[469]|11)-(?:0[1-9]|[12]\\d|30)|(?:02)-(?:0[1-9]|1\\d|2[0-8])))";
var aJ = new RegExp(`^${tJ}$`);
function sJ(Q) {
  return typeof Q.precision === "number" ? Q.precision === -1 ? "(?:[01]\\d|2[0-3]):[0-5]\\d" : Q.precision === 0 ? "(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d" : `(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d\\.\\d{${Q.precision}}` : "(?:[01]\\d|2[0-3]):[0-5]\\d(?::[0-5]\\d(?:\\.\\d+)?)?";
}
function eJ(Q) {
  return new RegExp(`^${sJ(Q)}$`);
}
function QG(Q) {
  let $ = sJ({ precision: Q.precision }), X = ["Z"];
  if (Q.local) X.push("");
  if (Q.offset) X.push("([+-]\\d{2}:\\d{2})");
  let Y = `${$}(?:${X.join("|")})`;
  return new RegExp(`^${tJ}T(?:${Y})$`);
}
var $G = (Q) => {
  let $ = Q ? `[\\s\\S]{${Q?.minimum ?? 0},${Q?.maximum ?? ""}}` : "[\\s\\S]*";
  return new RegExp(`^${$}$`);
};
var XG = /^\d+$/;
var YG = /^-?\d+(?:\.\d+)?/i;
var JG = /true|false/i;
var GG = /null/i;
var WG = /^[^A-Z]*$/;
var HG = /^[^a-z]*$/;
var j0 = O("$ZodCheck", (Q, $) => {
  var X;
  Q._zod ?? (Q._zod = {}), Q._zod.def = $, (X = Q._zod).onattach ?? (X.onattach = []);
});
var BG = { number: "number", bigint: "bigint", object: "date" };
var C$ = O("$ZodCheckLessThan", (Q, $) => {
  j0.init(Q, $);
  let X = BG[typeof $.value];
  Q._zod.onattach.push((Y) => {
    let J = Y._zod.bag, G = ($.inclusive ? J.maximum : J.exclusiveMaximum) ?? Number.POSITIVE_INFINITY;
    if ($.value < G) if ($.inclusive) J.maximum = $.value;
    else J.exclusiveMaximum = $.value;
  }), Q._zod.check = (Y) => {
    if ($.inclusive ? Y.value <= $.value : Y.value < $.value) return;
    Y.issues.push({ origin: X, code: "too_big", maximum: $.value, input: Y.value, inclusive: $.inclusive, inst: Q, continue: !$.abort });
  };
});
var S$ = O("$ZodCheckGreaterThan", (Q, $) => {
  j0.init(Q, $);
  let X = BG[typeof $.value];
  Q._zod.onattach.push((Y) => {
    let J = Y._zod.bag, G = ($.inclusive ? J.minimum : J.exclusiveMinimum) ?? Number.NEGATIVE_INFINITY;
    if ($.value > G) if ($.inclusive) J.minimum = $.value;
    else J.exclusiveMinimum = $.value;
  }), Q._zod.check = (Y) => {
    if ($.inclusive ? Y.value >= $.value : Y.value > $.value) return;
    Y.issues.push({ origin: X, code: "too_small", minimum: $.value, input: Y.value, inclusive: $.inclusive, inst: Q, continue: !$.abort });
  };
});
var zG = O("$ZodCheckMultipleOf", (Q, $) => {
  j0.init(Q, $), Q._zod.onattach.push((X) => {
    var Y;
    (Y = X._zod.bag).multipleOf ?? (Y.multipleOf = $.value);
  }), Q._zod.check = (X) => {
    if (typeof X.value !== typeof $.value) throw Error("Cannot mix number and bigint in multiple_of check.");
    if (typeof X.value === "bigint" ? X.value % $.value === BigInt(0) : q$(X.value, $.value) === 0) return;
    X.issues.push({ origin: typeof X.value, code: "not_multiple_of", divisor: $.value, input: X.value, inst: Q, continue: !$.abort });
  };
});
var KG = O("$ZodCheckNumberFormat", (Q, $) => {
  j0.init(Q, $), $.format = $.format || "float64";
  let X = $.format?.includes("int"), Y = X ? "int" : "number", [J, G] = O$[$.format];
  Q._zod.onattach.push((W) => {
    let H = W._zod.bag;
    if (H.format = $.format, H.minimum = J, H.maximum = G, X) H.pattern = XG;
  }), Q._zod.check = (W) => {
    let H = W.value;
    if (X) {
      if (!Number.isInteger(H)) {
        W.issues.push({ expected: Y, format: $.format, code: "invalid_type", input: H, inst: Q });
        return;
      }
      if (!Number.isSafeInteger(H)) {
        if (H > 0) W.issues.push({ input: H, code: "too_big", maximum: Number.MAX_SAFE_INTEGER, note: "Integers must be within the safe integer range.", inst: Q, origin: Y, continue: !$.abort });
        else W.issues.push({ input: H, code: "too_small", minimum: Number.MIN_SAFE_INTEGER, note: "Integers must be within the safe integer range.", inst: Q, origin: Y, continue: !$.abort });
        return;
      }
    }
    if (H < J) W.issues.push({ origin: "number", input: H, code: "too_small", minimum: J, inclusive: true, inst: Q, continue: !$.abort });
    if (H > G) W.issues.push({ origin: "number", input: H, code: "too_big", maximum: G, inst: Q });
  };
});
var VG = O("$ZodCheckMaxLength", (Q, $) => {
  j0.init(Q, $), Q._zod.when = (X) => {
    let Y = X.value;
    return !s4(Y) && Y.length !== void 0;
  }, Q._zod.onattach.push((X) => {
    let Y = X._zod.bag.maximum ?? Number.POSITIVE_INFINITY;
    if ($.maximum < Y) X._zod.bag.maximum = $.maximum;
  }), Q._zod.check = (X) => {
    let Y = X.value;
    if (Y.length <= $.maximum) return;
    let G = Q9(Y);
    X.issues.push({ origin: G, code: "too_big", maximum: $.maximum, inclusive: true, input: Y, inst: Q, continue: !$.abort });
  };
});
var qG = O("$ZodCheckMinLength", (Q, $) => {
  j0.init(Q, $), Q._zod.when = (X) => {
    let Y = X.value;
    return !s4(Y) && Y.length !== void 0;
  }, Q._zod.onattach.push((X) => {
    let Y = X._zod.bag.minimum ?? Number.NEGATIVE_INFINITY;
    if ($.minimum > Y) X._zod.bag.minimum = $.minimum;
  }), Q._zod.check = (X) => {
    let Y = X.value;
    if (Y.length >= $.minimum) return;
    let G = Q9(Y);
    X.issues.push({ origin: G, code: "too_small", minimum: $.minimum, inclusive: true, input: Y, inst: Q, continue: !$.abort });
  };
});
var UG = O("$ZodCheckLengthEquals", (Q, $) => {
  j0.init(Q, $), Q._zod.when = (X) => {
    let Y = X.value;
    return !s4(Y) && Y.length !== void 0;
  }, Q._zod.onattach.push((X) => {
    let Y = X._zod.bag;
    Y.minimum = $.length, Y.maximum = $.length, Y.length = $.length;
  }), Q._zod.check = (X) => {
    let Y = X.value, J = Y.length;
    if (J === $.length) return;
    let G = Q9(Y), W = J > $.length;
    X.issues.push({ origin: G, ...W ? { code: "too_big", maximum: $.length } : { code: "too_small", minimum: $.length }, inclusive: true, exact: true, input: X.value, inst: Q, continue: !$.abort });
  };
});
var X9 = O("$ZodCheckStringFormat", (Q, $) => {
  var X, Y;
  if (j0.init(Q, $), Q._zod.onattach.push((J) => {
    let G = J._zod.bag;
    if (G.format = $.format, $.pattern) G.patterns ?? (G.patterns = /* @__PURE__ */ new Set()), G.patterns.add($.pattern);
  }), $.pattern) (X = Q._zod).check ?? (X.check = (J) => {
    if ($.pattern.lastIndex = 0, $.pattern.test(J.value)) return;
    J.issues.push({ origin: "string", code: "invalid_format", format: $.format, input: J.value, ...$.pattern ? { pattern: $.pattern.toString() } : {}, inst: Q, continue: !$.abort });
  });
  else (Y = Q._zod).check ?? (Y.check = () => {
  });
});
var LG = O("$ZodCheckRegex", (Q, $) => {
  X9.init(Q, $), Q._zod.check = (X) => {
    if ($.pattern.lastIndex = 0, $.pattern.test(X.value)) return;
    X.issues.push({ origin: "string", code: "invalid_format", format: "regex", input: X.value, pattern: $.pattern.toString(), inst: Q, continue: !$.abort });
  };
});
var FG = O("$ZodCheckLowerCase", (Q, $) => {
  $.pattern ?? ($.pattern = WG), X9.init(Q, $);
});
var NG = O("$ZodCheckUpperCase", (Q, $) => {
  $.pattern ?? ($.pattern = HG), X9.init(Q, $);
});
var OG = O("$ZodCheckIncludes", (Q, $) => {
  j0.init(Q, $);
  let X = o1($.includes), Y = new RegExp(typeof $.position === "number" ? `^.{${$.position}}${X}` : X);
  $.pattern = Y, Q._zod.onattach.push((J) => {
    let G = J._zod.bag;
    G.patterns ?? (G.patterns = /* @__PURE__ */ new Set()), G.patterns.add(Y);
  }), Q._zod.check = (J) => {
    if (J.value.includes($.includes, $.position)) return;
    J.issues.push({ origin: "string", code: "invalid_format", format: "includes", includes: $.includes, input: J.value, inst: Q, continue: !$.abort });
  };
});
var DG = O("$ZodCheckStartsWith", (Q, $) => {
  j0.init(Q, $);
  let X = new RegExp(`^${o1($.prefix)}.*`);
  $.pattern ?? ($.pattern = X), Q._zod.onattach.push((Y) => {
    let J = Y._zod.bag;
    J.patterns ?? (J.patterns = /* @__PURE__ */ new Set()), J.patterns.add(X);
  }), Q._zod.check = (Y) => {
    if (Y.value.startsWith($.prefix)) return;
    Y.issues.push({ origin: "string", code: "invalid_format", format: "starts_with", prefix: $.prefix, input: Y.value, inst: Q, continue: !$.abort });
  };
});
var wG = O("$ZodCheckEndsWith", (Q, $) => {
  j0.init(Q, $);
  let X = new RegExp(`.*${o1($.suffix)}$`);
  $.pattern ?? ($.pattern = X), Q._zod.onattach.push((Y) => {
    let J = Y._zod.bag;
    J.patterns ?? (J.patterns = /* @__PURE__ */ new Set()), J.patterns.add(X);
  }), Q._zod.check = (Y) => {
    if (Y.value.endsWith($.suffix)) return;
    Y.issues.push({ origin: "string", code: "invalid_format", format: "ends_with", suffix: $.suffix, input: Y.value, inst: Q, continue: !$.abort });
  };
});
var MG = O("$ZodCheckOverwrite", (Q, $) => {
  j0.init(Q, $), Q._zod.check = (X) => {
    X.value = $.tx(X.value);
  };
});
var _$ = class {
  constructor(Q = []) {
    if (this.content = [], this.indent = 0, this) this.args = Q;
  }
  indented(Q) {
    this.indent += 1, Q(this), this.indent -= 1;
  }
  write(Q) {
    if (typeof Q === "function") {
      Q(this, { execution: "sync" }), Q(this, { execution: "async" });
      return;
    }
    let X = Q.split(`
`).filter((G) => G), Y = Math.min(...X.map((G) => G.length - G.trimStart().length)), J = X.map((G) => G.slice(Y)).map((G) => " ".repeat(this.indent * 2) + G);
    for (let G of J) this.content.push(G);
  }
  compile() {
    let Q = Function, $ = this?.args, Y = [...(this?.content ?? [""]).map((J) => `  ${J}`)];
    return new Q(...$, Y.join(`
`));
  }
};
var jG = { major: 4, minor: 0, patch: 0 };
var e = O("$ZodType", (Q, $) => {
  var X;
  Q ?? (Q = {}), Q._zod.def = $, Q._zod.bag = Q._zod.bag || {}, Q._zod.version = jG;
  let Y = [...Q._zod.def.checks ?? []];
  if (Q._zod.traits.has("$ZodCheck")) Y.unshift(Q);
  for (let J of Y) for (let G of J._zod.onattach) G(Q);
  if (Y.length === 0) (X = Q._zod).deferred ?? (X.deferred = []), Q._zod.deferred?.push(() => {
    Q._zod.run = Q._zod.parse;
  });
  else {
    let J = (G, W, H) => {
      let B = A6(G), z;
      for (let K of W) {
        if (K._zod.when) {
          if (!K._zod.when(G)) continue;
        } else if (B) continue;
        let U = G.issues.length, q = K._zod.check(G);
        if (q instanceof Promise && H?.async === false) throw new n1();
        if (z || q instanceof Promise) z = (z ?? Promise.resolve()).then(async () => {
          if (await q, G.issues.length === U) return;
          if (!B) B = A6(G, U);
        });
        else {
          if (G.issues.length === U) continue;
          if (!B) B = A6(G, U);
        }
      }
      if (z) return z.then(() => {
        return G;
      });
      return G;
    };
    Q._zod.run = (G, W) => {
      let H = Q._zod.parse(G, W);
      if (H instanceof Promise) {
        if (W.async === false) throw new n1();
        return H.then((B) => J(B, Y, W));
      }
      return J(H, Y, W);
    };
  }
  Q["~standard"] = { validate: (J) => {
    try {
      let G = j6(Q, J);
      return G.success ? { value: G.data } : { issues: G.error?.issues };
    } catch (G) {
      return R6(Q, J).then((W) => W.success ? { value: W.data } : { issues: W.error?.issues });
    }
  }, vendor: "zod", version: 1 };
});
var Y9 = O("$ZodString", (Q, $) => {
  e.init(Q, $), Q._zod.pattern = [...Q?._zod.bag?.patterns ?? []].pop() ?? $G(Q._zod.bag), Q._zod.parse = (X, Y) => {
    if ($.coerce) try {
      X.value = String(X.value);
    } catch (J) {
    }
    if (typeof X.value === "string") return X;
    return X.issues.push({ expected: "string", code: "invalid_type", input: X.value, inst: Q }), X;
  };
});
var J0 = O("$ZodStringFormat", (Q, $) => {
  X9.init(Q, $), Y9.init(Q, $);
});
var v$ = O("$ZodGUID", (Q, $) => {
  $.pattern ?? ($.pattern = uJ), J0.init(Q, $);
});
var T$ = O("$ZodUUID", (Q, $) => {
  if ($.version) {
    let Y = { v1: 1, v2: 2, v3: 3, v4: 4, v5: 5, v6: 6, v7: 7, v8: 8 }[$.version];
    if (Y === void 0) throw Error(`Invalid UUID version: "${$.version}"`);
    $.pattern ?? ($.pattern = b$(Y));
  } else $.pattern ?? ($.pattern = b$());
  J0.init(Q, $);
});
var x$ = O("$ZodEmail", (Q, $) => {
  $.pattern ?? ($.pattern = mJ), J0.init(Q, $);
});
var y$ = O("$ZodURL", (Q, $) => {
  J0.init(Q, $), Q._zod.check = (X) => {
    try {
      let Y = X.value, J = new URL(Y), G = J.href;
      if ($.hostname) {
        if ($.hostname.lastIndex = 0, !$.hostname.test(J.hostname)) X.issues.push({ code: "invalid_format", format: "url", note: "Invalid hostname", pattern: oJ.source, input: X.value, inst: Q, continue: !$.abort });
      }
      if ($.protocol) {
        if ($.protocol.lastIndex = 0, !$.protocol.test(J.protocol.endsWith(":") ? J.protocol.slice(0, -1) : J.protocol)) X.issues.push({ code: "invalid_format", format: "url", note: "Invalid protocol", pattern: $.protocol.source, input: X.value, inst: Q, continue: !$.abort });
      }
      if (!Y.endsWith("/") && G.endsWith("/")) X.value = G.slice(0, -1);
      else X.value = G;
      return;
    } catch (Y) {
      X.issues.push({ code: "invalid_format", format: "url", input: X.value, inst: Q, continue: !$.abort });
    }
  };
});
var g$ = O("$ZodEmoji", (Q, $) => {
  $.pattern ?? ($.pattern = lJ()), J0.init(Q, $);
});
var h$ = O("$ZodNanoID", (Q, $) => {
  $.pattern ?? ($.pattern = hJ), J0.init(Q, $);
});
var f$ = O("$ZodCUID", (Q, $) => {
  $.pattern ?? ($.pattern = vJ), J0.init(Q, $);
});
var u$ = O("$ZodCUID2", (Q, $) => {
  $.pattern ?? ($.pattern = TJ), J0.init(Q, $);
});
var m$ = O("$ZodULID", (Q, $) => {
  $.pattern ?? ($.pattern = xJ), J0.init(Q, $);
});
var l$ = O("$ZodXID", (Q, $) => {
  $.pattern ?? ($.pattern = yJ), J0.init(Q, $);
});
var c$ = O("$ZodKSUID", (Q, $) => {
  $.pattern ?? ($.pattern = gJ), J0.init(Q, $);
});
var kG = O("$ZodISODateTime", (Q, $) => {
  $.pattern ?? ($.pattern = QG($)), J0.init(Q, $);
});
var vG = O("$ZodISODate", (Q, $) => {
  $.pattern ?? ($.pattern = aJ), J0.init(Q, $);
});
var TG = O("$ZodISOTime", (Q, $) => {
  $.pattern ?? ($.pattern = eJ($)), J0.init(Q, $);
});
var xG = O("$ZodISODuration", (Q, $) => {
  $.pattern ?? ($.pattern = fJ), J0.init(Q, $);
});
var p$ = O("$ZodIPv4", (Q, $) => {
  $.pattern ?? ($.pattern = cJ), J0.init(Q, $), Q._zod.onattach.push((X) => {
    let Y = X._zod.bag;
    Y.format = "ipv4";
  });
});
var d$ = O("$ZodIPv6", (Q, $) => {
  $.pattern ?? ($.pattern = pJ), J0.init(Q, $), Q._zod.onattach.push((X) => {
    let Y = X._zod.bag;
    Y.format = "ipv6";
  }), Q._zod.check = (X) => {
    try {
      new URL(`http://[${X.value}]`);
    } catch {
      X.issues.push({ code: "invalid_format", format: "ipv6", input: X.value, inst: Q, continue: !$.abort });
    }
  };
});
var i$ = O("$ZodCIDRv4", (Q, $) => {
  $.pattern ?? ($.pattern = dJ), J0.init(Q, $);
});
var n$ = O("$ZodCIDRv6", (Q, $) => {
  $.pattern ?? ($.pattern = iJ), J0.init(Q, $), Q._zod.check = (X) => {
    let [Y, J] = X.value.split("/");
    try {
      if (!J) throw Error();
      let G = Number(J);
      if (`${G}` !== J) throw Error();
      if (G < 0 || G > 128) throw Error();
      new URL(`http://[${Y}]`);
    } catch {
      X.issues.push({ code: "invalid_format", format: "cidrv6", input: X.value, inst: Q, continue: !$.abort });
    }
  };
});
function yG(Q) {
  if (Q === "") return true;
  if (Q.length % 4 !== 0) return false;
  try {
    return atob(Q), true;
  } catch {
    return false;
  }
}
var o$ = O("$ZodBase64", (Q, $) => {
  $.pattern ?? ($.pattern = nJ), J0.init(Q, $), Q._zod.onattach.push((X) => {
    X._zod.bag.contentEncoding = "base64";
  }), Q._zod.check = (X) => {
    if (yG(X.value)) return;
    X.issues.push({ code: "invalid_format", format: "base64", input: X.value, inst: Q, continue: !$.abort });
  };
});
function $N(Q) {
  if (!Z$.test(Q)) return false;
  let $ = Q.replace(/[-_]/g, (Y) => Y === "-" ? "+" : "/"), X = $.padEnd(Math.ceil($.length / 4) * 4, "=");
  return yG(X);
}
var r$ = O("$ZodBase64URL", (Q, $) => {
  $.pattern ?? ($.pattern = Z$), J0.init(Q, $), Q._zod.onattach.push((X) => {
    X._zod.bag.contentEncoding = "base64url";
  }), Q._zod.check = (X) => {
    if ($N(X.value)) return;
    X.issues.push({ code: "invalid_format", format: "base64url", input: X.value, inst: Q, continue: !$.abort });
  };
});
var t$ = O("$ZodE164", (Q, $) => {
  $.pattern ?? ($.pattern = rJ), J0.init(Q, $);
});
function XN(Q, $ = null) {
  try {
    let X = Q.split(".");
    if (X.length !== 3) return false;
    let [Y] = X;
    if (!Y) return false;
    let J = JSON.parse(atob(Y));
    if ("typ" in J && J?.typ !== "JWT") return false;
    if (!J.alg) return false;
    if ($ && (!("alg" in J) || J.alg !== $)) return false;
    return true;
  } catch {
    return false;
  }
}
var a$ = O("$ZodJWT", (Q, $) => {
  J0.init(Q, $), Q._zod.check = (X) => {
    if (XN(X.value, $.alg)) return;
    X.issues.push({ code: "invalid_format", format: "jwt", input: X.value, inst: Q, continue: !$.abort });
  };
});
var L8 = O("$ZodNumber", (Q, $) => {
  e.init(Q, $), Q._zod.pattern = Q._zod.bag.pattern ?? YG, Q._zod.parse = (X, Y) => {
    if ($.coerce) try {
      X.value = Number(X.value);
    } catch (W) {
    }
    let J = X.value;
    if (typeof J === "number" && !Number.isNaN(J) && Number.isFinite(J)) return X;
    let G = typeof J === "number" ? Number.isNaN(J) ? "NaN" : !Number.isFinite(J) ? "Infinity" : void 0 : void 0;
    return X.issues.push({ expected: "number", code: "invalid_type", input: J, inst: Q, ...G ? { received: G } : {} }), X;
  };
});
var s$ = O("$ZodNumber", (Q, $) => {
  KG.init(Q, $), L8.init(Q, $);
});
var e$ = O("$ZodBoolean", (Q, $) => {
  e.init(Q, $), Q._zod.pattern = JG, Q._zod.parse = (X, Y) => {
    if ($.coerce) try {
      X.value = Boolean(X.value);
    } catch (G) {
    }
    let J = X.value;
    if (typeof J === "boolean") return X;
    return X.issues.push({ expected: "boolean", code: "invalid_type", input: J, inst: Q }), X;
  };
});
var QX = O("$ZodNull", (Q, $) => {
  e.init(Q, $), Q._zod.pattern = GG, Q._zod.values = /* @__PURE__ */ new Set([null]), Q._zod.parse = (X, Y) => {
    let J = X.value;
    if (J === null) return X;
    return X.issues.push({ expected: "null", code: "invalid_type", input: J, inst: Q }), X;
  };
});
var $X = O("$ZodUnknown", (Q, $) => {
  e.init(Q, $), Q._zod.parse = (X) => X;
});
var XX = O("$ZodNever", (Q, $) => {
  e.init(Q, $), Q._zod.parse = (X, Y) => {
    return X.issues.push({ expected: "never", code: "invalid_type", input: X.value, inst: Q }), X;
  };
});
function RG(Q, $, X) {
  if (Q.issues.length) $.issues.push(...j1(X, Q.issues));
  $.value[X] = Q.value;
}
var YX = O("$ZodArray", (Q, $) => {
  e.init(Q, $), Q._zod.parse = (X, Y) => {
    let J = X.value;
    if (!Array.isArray(J)) return X.issues.push({ expected: "array", code: "invalid_type", input: J, inst: Q }), X;
    X.value = Array(J.length);
    let G = [];
    for (let W = 0; W < J.length; W++) {
      let H = J[W], B = $.element._zod.run({ value: H, issues: [] }, Y);
      if (B instanceof Promise) G.push(B.then((z) => RG(z, X, W)));
      else RG(B, X, W);
    }
    if (G.length) return Promise.all(G).then(() => X);
    return X;
  };
});
function U8(Q, $, X) {
  if (Q.issues.length) $.issues.push(...j1(X, Q.issues));
  $.value[X] = Q.value;
}
function IG(Q, $, X, Y) {
  if (Q.issues.length) if (Y[X] === void 0) if (X in Y) $.value[X] = void 0;
  else $.value[X] = Q.value;
  else $.issues.push(...j1(X, Q.issues));
  else if (Q.value === void 0) {
    if (X in Y) $.value[X] = void 0;
  } else $.value[X] = Q.value;
}
var F8 = O("$ZodObject", (Q, $) => {
  e.init(Q, $);
  let X = a4(() => {
    let U = Object.keys($.shape);
    for (let V of U) if (!($.shape[V] instanceof e)) throw Error(`Invalid element at key "${V}": expected a Zod schema`);
    let q = N$($.shape);
    return { shape: $.shape, keys: U, keySet: new Set(U), numKeys: U.length, optionalKeys: new Set(q) };
  });
  Y0(Q._zod, "propValues", () => {
    let U = $.shape, q = {};
    for (let V in U) {
      let L = U[V]._zod;
      if (L.values) {
        q[V] ?? (q[V] = /* @__PURE__ */ new Set());
        for (let F of L.values) q[V].add(F);
      }
    }
    return q;
  });
  let Y = (U) => {
    let q = new _$(["shape", "payload", "ctx"]), V = X.value, L = (M) => {
      let R = M6(M);
      return `shape[${R}]._zod.run({ value: input[${R}], issues: [] }, ctx)`;
    };
    q.write("const input = payload.value;");
    let F = /* @__PURE__ */ Object.create(null), w = 0;
    for (let M of V.keys) F[M] = `key_${w++}`;
    q.write("const newResult = {}");
    for (let M of V.keys) if (V.optionalKeys.has(M)) {
      let R = F[M];
      q.write(`const ${R} = ${L(M)};`);
      let Z = M6(M);
      q.write(`
        if (${R}.issues.length) {
          if (input[${Z}] === undefined) {
            if (${Z} in input) {
              newResult[${Z}] = undefined;
            }
          } else {
            payload.issues = payload.issues.concat(
              ${R}.issues.map((iss) => ({
                ...iss,
                path: iss.path ? [${Z}, ...iss.path] : [${Z}],
              }))
            );
          }
        } else if (${R}.value === undefined) {
          if (${Z} in input) newResult[${Z}] = undefined;
        } else {
          newResult[${Z}] = ${R}.value;
        }
        `);
    } else {
      let R = F[M];
      q.write(`const ${R} = ${L(M)};`), q.write(`
          if (${R}.issues.length) payload.issues = payload.issues.concat(${R}.issues.map(iss => ({
            ...iss,
            path: iss.path ? [${M6(M)}, ...iss.path] : [${M6(M)}]
          })));`), q.write(`newResult[${M6(M)}] = ${R}.value`);
    }
    q.write("payload.value = newResult;"), q.write("return payload;");
    let D = q.compile();
    return (M, R) => D(U, M, R);
  }, J, G = K4, W = !H8.jitless, B = W && L$.value, z = $.catchall, K;
  Q._zod.parse = (U, q) => {
    K ?? (K = X.value);
    let V = U.value;
    if (!G(V)) return U.issues.push({ expected: "object", code: "invalid_type", input: V, inst: Q }), U;
    let L = [];
    if (W && B && q?.async === false && q.jitless !== true) {
      if (!J) J = Y($.shape);
      U = J(U, q);
    } else {
      U.value = {};
      let R = K.shape;
      for (let Z of K.keys) {
        let v = R[Z], O0 = v._zod.run({ value: V[Z], issues: [] }, q), D0 = v._zod.optin === "optional" && v._zod.optout === "optional";
        if (O0 instanceof Promise) L.push(O0.then((d0) => D0 ? IG(d0, U, Z, V) : U8(d0, U, Z)));
        else if (D0) IG(O0, U, Z, V);
        else U8(O0, U, Z);
      }
    }
    if (!z) return L.length ? Promise.all(L).then(() => U) : U;
    let F = [], w = K.keySet, D = z._zod, M = D.def.type;
    for (let R of Object.keys(V)) {
      if (w.has(R)) continue;
      if (M === "never") {
        F.push(R);
        continue;
      }
      let Z = D.run({ value: V[R], issues: [] }, q);
      if (Z instanceof Promise) L.push(Z.then((v) => U8(v, U, R)));
      else U8(Z, U, R);
    }
    if (F.length) U.issues.push({ code: "unrecognized_keys", keys: F, input: V, inst: Q });
    if (!L.length) return U;
    return Promise.all(L).then(() => {
      return U;
    });
  };
});
function EG(Q, $, X, Y) {
  for (let J of Q) if (J.issues.length === 0) return $.value = J.value, $;
  return $.issues.push({ code: "invalid_union", input: $.value, inst: X, errors: Q.map((J) => J.issues.map((G) => W1(G, Y, l0()))) }), $;
}
var N8 = O("$ZodUnion", (Q, $) => {
  e.init(Q, $), Y0(Q._zod, "optin", () => $.options.some((X) => X._zod.optin === "optional") ? "optional" : void 0), Y0(Q._zod, "optout", () => $.options.some((X) => X._zod.optout === "optional") ? "optional" : void 0), Y0(Q._zod, "values", () => {
    if ($.options.every((X) => X._zod.values)) return new Set($.options.flatMap((X) => Array.from(X._zod.values)));
    return;
  }), Y0(Q._zod, "pattern", () => {
    if ($.options.every((X) => X._zod.pattern)) {
      let X = $.options.map((Y) => Y._zod.pattern);
      return new RegExp(`^(${X.map((Y) => e4(Y.source)).join("|")})$`);
    }
    return;
  }), Q._zod.parse = (X, Y) => {
    let J = false, G = [];
    for (let W of $.options) {
      let H = W._zod.run({ value: X.value, issues: [] }, Y);
      if (H instanceof Promise) G.push(H), J = true;
      else {
        if (H.issues.length === 0) return H;
        G.push(H);
      }
    }
    if (!J) return EG(G, X, Q, Y);
    return Promise.all(G).then((W) => {
      return EG(W, X, Q, Y);
    });
  };
});
var JX = O("$ZodDiscriminatedUnion", (Q, $) => {
  N8.init(Q, $);
  let X = Q._zod.parse;
  Y0(Q._zod, "propValues", () => {
    let J = {};
    for (let G of $.options) {
      let W = G._zod.propValues;
      if (!W || Object.keys(W).length === 0) throw Error(`Invalid discriminated union option at index "${$.options.indexOf(G)}"`);
      for (let [H, B] of Object.entries(W)) {
        if (!J[H]) J[H] = /* @__PURE__ */ new Set();
        for (let z of B) J[H].add(z);
      }
    }
    return J;
  });
  let Y = a4(() => {
    let J = $.options, G = /* @__PURE__ */ new Map();
    for (let W of J) {
      let H = W._zod.propValues[$.discriminator];
      if (!H || H.size === 0) throw Error(`Invalid discriminated union option at index "${$.options.indexOf(W)}"`);
      for (let B of H) {
        if (G.has(B)) throw Error(`Duplicate discriminator value "${String(B)}"`);
        G.set(B, W);
      }
    }
    return G;
  });
  Q._zod.parse = (J, G) => {
    let W = J.value;
    if (!K4(W)) return J.issues.push({ code: "invalid_type", expected: "object", input: W, inst: Q }), J;
    let H = Y.value.get(W?.[$.discriminator]);
    if (H) return H._zod.run(J, G);
    if ($.unionFallback) return X(J, G);
    return J.issues.push({ code: "invalid_union", errors: [], note: "No matching discriminator", input: W, path: [$.discriminator], inst: Q }), J;
  };
});
var GX = O("$ZodIntersection", (Q, $) => {
  e.init(Q, $), Q._zod.parse = (X, Y) => {
    let J = X.value, G = $.left._zod.run({ value: J, issues: [] }, Y), W = $.right._zod.run({ value: J, issues: [] }, Y);
    if (G instanceof Promise || W instanceof Promise) return Promise.all([G, W]).then(([B, z]) => {
      return PG(X, B, z);
    });
    return PG(X, G, W);
  };
});
function k$(Q, $) {
  if (Q === $) return { valid: true, data: Q };
  if (Q instanceof Date && $ instanceof Date && +Q === +$) return { valid: true, data: Q };
  if (V4(Q) && V4($)) {
    let X = Object.keys($), Y = Object.keys(Q).filter((G) => X.indexOf(G) !== -1), J = { ...Q, ...$ };
    for (let G of Y) {
      let W = k$(Q[G], $[G]);
      if (!W.valid) return { valid: false, mergeErrorPath: [G, ...W.mergeErrorPath] };
      J[G] = W.data;
    }
    return { valid: true, data: J };
  }
  if (Array.isArray(Q) && Array.isArray($)) {
    if (Q.length !== $.length) return { valid: false, mergeErrorPath: [] };
    let X = [];
    for (let Y = 0; Y < Q.length; Y++) {
      let J = Q[Y], G = $[Y], W = k$(J, G);
      if (!W.valid) return { valid: false, mergeErrorPath: [Y, ...W.mergeErrorPath] };
      X.push(W.data);
    }
    return { valid: true, data: X };
  }
  return { valid: false, mergeErrorPath: [] };
}
function PG(Q, $, X) {
  if ($.issues.length) Q.issues.push(...$.issues);
  if (X.issues.length) Q.issues.push(...X.issues);
  if (A6(Q)) return Q;
  let Y = k$($.value, X.value);
  if (!Y.valid) throw Error(`Unmergable intersection. Error path: ${JSON.stringify(Y.mergeErrorPath)}`);
  return Q.value = Y.data, Q;
}
var WX = O("$ZodRecord", (Q, $) => {
  e.init(Q, $), Q._zod.parse = (X, Y) => {
    let J = X.value;
    if (!V4(J)) return X.issues.push({ expected: "record", code: "invalid_type", input: J, inst: Q }), X;
    let G = [];
    if ($.keyType._zod.values) {
      let W = $.keyType._zod.values;
      X.value = {};
      for (let B of W) if (typeof B === "string" || typeof B === "number" || typeof B === "symbol") {
        let z = $.valueType._zod.run({ value: J[B], issues: [] }, Y);
        if (z instanceof Promise) G.push(z.then((K) => {
          if (K.issues.length) X.issues.push(...j1(B, K.issues));
          X.value[B] = K.value;
        }));
        else {
          if (z.issues.length) X.issues.push(...j1(B, z.issues));
          X.value[B] = z.value;
        }
      }
      let H;
      for (let B in J) if (!W.has(B)) H = H ?? [], H.push(B);
      if (H && H.length > 0) X.issues.push({ code: "unrecognized_keys", input: J, inst: Q, keys: H });
    } else {
      X.value = {};
      for (let W of Reflect.ownKeys(J)) {
        if (W === "__proto__") continue;
        let H = $.keyType._zod.run({ value: W, issues: [] }, Y);
        if (H instanceof Promise) throw Error("Async schemas not supported in object keys currently");
        if (H.issues.length) {
          X.issues.push({ origin: "record", code: "invalid_key", issues: H.issues.map((z) => W1(z, Y, l0())), input: W, path: [W], inst: Q }), X.value[H.value] = H.value;
          continue;
        }
        let B = $.valueType._zod.run({ value: J[W], issues: [] }, Y);
        if (B instanceof Promise) G.push(B.then((z) => {
          if (z.issues.length) X.issues.push(...j1(W, z.issues));
          X.value[H.value] = z.value;
        }));
        else {
          if (B.issues.length) X.issues.push(...j1(W, B.issues));
          X.value[H.value] = B.value;
        }
      }
    }
    if (G.length) return Promise.all(G).then(() => X);
    return X;
  };
});
var HX = O("$ZodEnum", (Q, $) => {
  e.init(Q, $);
  let X = t4($.entries);
  Q._zod.values = new Set(X), Q._zod.pattern = new RegExp(`^(${X.filter((Y) => F$.has(typeof Y)).map((Y) => typeof Y === "string" ? o1(Y) : Y.toString()).join("|")})$`), Q._zod.parse = (Y, J) => {
    let G = Y.value;
    if (Q._zod.values.has(G)) return Y;
    return Y.issues.push({ code: "invalid_value", values: X, input: G, inst: Q }), Y;
  };
});
var BX = O("$ZodLiteral", (Q, $) => {
  e.init(Q, $), Q._zod.values = new Set($.values), Q._zod.pattern = new RegExp(`^(${$.values.map((X) => typeof X === "string" ? o1(X) : X ? X.toString() : String(X)).join("|")})$`), Q._zod.parse = (X, Y) => {
    let J = X.value;
    if (Q._zod.values.has(J)) return X;
    return X.issues.push({ code: "invalid_value", values: $.values, input: J, inst: Q }), X;
  };
});
var zX = O("$ZodTransform", (Q, $) => {
  e.init(Q, $), Q._zod.parse = (X, Y) => {
    let J = $.transform(X.value, X);
    if (Y.async) return (J instanceof Promise ? J : Promise.resolve(J)).then((W) => {
      return X.value = W, X;
    });
    if (J instanceof Promise) throw new n1();
    return X.value = J, X;
  };
});
var KX = O("$ZodOptional", (Q, $) => {
  e.init(Q, $), Q._zod.optin = "optional", Q._zod.optout = "optional", Y0(Q._zod, "values", () => {
    return $.innerType._zod.values ? /* @__PURE__ */ new Set([...$.innerType._zod.values, void 0]) : void 0;
  }), Y0(Q._zod, "pattern", () => {
    let X = $.innerType._zod.pattern;
    return X ? new RegExp(`^(${e4(X.source)})?$`) : void 0;
  }), Q._zod.parse = (X, Y) => {
    if ($.innerType._zod.optin === "optional") return $.innerType._zod.run(X, Y);
    if (X.value === void 0) return X;
    return $.innerType._zod.run(X, Y);
  };
});
var VX = O("$ZodNullable", (Q, $) => {
  e.init(Q, $), Y0(Q._zod, "optin", () => $.innerType._zod.optin), Y0(Q._zod, "optout", () => $.innerType._zod.optout), Y0(Q._zod, "pattern", () => {
    let X = $.innerType._zod.pattern;
    return X ? new RegExp(`^(${e4(X.source)}|null)$`) : void 0;
  }), Y0(Q._zod, "values", () => {
    return $.innerType._zod.values ? /* @__PURE__ */ new Set([...$.innerType._zod.values, null]) : void 0;
  }), Q._zod.parse = (X, Y) => {
    if (X.value === null) return X;
    return $.innerType._zod.run(X, Y);
  };
});
var qX = O("$ZodDefault", (Q, $) => {
  e.init(Q, $), Q._zod.optin = "optional", Y0(Q._zod, "values", () => $.innerType._zod.values), Q._zod.parse = (X, Y) => {
    if (X.value === void 0) return X.value = $.defaultValue, X;
    let J = $.innerType._zod.run(X, Y);
    if (J instanceof Promise) return J.then((G) => bG(G, $));
    return bG(J, $);
  };
});
function bG(Q, $) {
  if (Q.value === void 0) Q.value = $.defaultValue;
  return Q;
}
var UX = O("$ZodPrefault", (Q, $) => {
  e.init(Q, $), Q._zod.optin = "optional", Y0(Q._zod, "values", () => $.innerType._zod.values), Q._zod.parse = (X, Y) => {
    if (X.value === void 0) X.value = $.defaultValue;
    return $.innerType._zod.run(X, Y);
  };
});
var LX = O("$ZodNonOptional", (Q, $) => {
  e.init(Q, $), Y0(Q._zod, "values", () => {
    let X = $.innerType._zod.values;
    return X ? new Set([...X].filter((Y) => Y !== void 0)) : void 0;
  }), Q._zod.parse = (X, Y) => {
    let J = $.innerType._zod.run(X, Y);
    if (J instanceof Promise) return J.then((G) => ZG(G, Q));
    return ZG(J, Q);
  };
});
function ZG(Q, $) {
  if (!Q.issues.length && Q.value === void 0) Q.issues.push({ code: "invalid_type", expected: "nonoptional", input: Q.value, inst: $ });
  return Q;
}
var FX = O("$ZodCatch", (Q, $) => {
  e.init(Q, $), Q._zod.optin = "optional", Y0(Q._zod, "optout", () => $.innerType._zod.optout), Y0(Q._zod, "values", () => $.innerType._zod.values), Q._zod.parse = (X, Y) => {
    let J = $.innerType._zod.run(X, Y);
    if (J instanceof Promise) return J.then((G) => {
      if (X.value = G.value, G.issues.length) X.value = $.catchValue({ ...X, error: { issues: G.issues.map((W) => W1(W, Y, l0())) }, input: X.value }), X.issues = [];
      return X;
    });
    if (X.value = J.value, J.issues.length) X.value = $.catchValue({ ...X, error: { issues: J.issues.map((G) => W1(G, Y, l0())) }, input: X.value }), X.issues = [];
    return X;
  };
});
var NX = O("$ZodPipe", (Q, $) => {
  e.init(Q, $), Y0(Q._zod, "values", () => $.in._zod.values), Y0(Q._zod, "optin", () => $.in._zod.optin), Y0(Q._zod, "optout", () => $.out._zod.optout), Q._zod.parse = (X, Y) => {
    let J = $.in._zod.run(X, Y);
    if (J instanceof Promise) return J.then((G) => CG(G, $, Y));
    return CG(J, $, Y);
  };
});
function CG(Q, $, X) {
  if (A6(Q)) return Q;
  return $.out._zod.run({ value: Q.value, issues: Q.issues }, X);
}
var OX = O("$ZodReadonly", (Q, $) => {
  e.init(Q, $), Y0(Q._zod, "propValues", () => $.innerType._zod.propValues), Y0(Q._zod, "values", () => $.innerType._zod.values), Y0(Q._zod, "optin", () => $.innerType._zod.optin), Y0(Q._zod, "optout", () => $.innerType._zod.optout), Q._zod.parse = (X, Y) => {
    let J = $.innerType._zod.run(X, Y);
    if (J instanceof Promise) return J.then(SG);
    return SG(J);
  };
});
function SG(Q) {
  return Q.value = Object.freeze(Q.value), Q;
}
var DX = O("$ZodCustom", (Q, $) => {
  j0.init(Q, $), e.init(Q, $), Q._zod.parse = (X, Y) => {
    return X;
  }, Q._zod.check = (X) => {
    let Y = X.value, J = $.fn(Y);
    if (J instanceof Promise) return J.then((G) => _G(G, X, Y, Q));
    _G(J, X, Y, Q);
    return;
  };
});
function _G(Q, $, X, Y) {
  if (!Q) {
    let J = { code: "custom", input: X, inst: Y, path: [...Y._zod.def.path ?? []], continue: !Y._zod.def.abort };
    if (Y._zod.def.params) J.params = Y._zod.def.params;
    $.issues.push(D$(J));
  }
}
var YN = (Q) => {
  let $ = typeof Q;
  switch ($) {
    case "number":
      return Number.isNaN(Q) ? "NaN" : "number";
    case "object": {
      if (Array.isArray(Q)) return "array";
      if (Q === null) return "null";
      if (Object.getPrototypeOf(Q) !== Object.prototype && Q.constructor) return Q.constructor.name;
    }
  }
  return $;
};
var JN = () => {
  let Q = { string: { unit: "characters", verb: "to have" }, file: { unit: "bytes", verb: "to have" }, array: { unit: "items", verb: "to have" }, set: { unit: "items", verb: "to have" } };
  function $(Y) {
    return Q[Y] ?? null;
  }
  let X = { regex: "input", email: "email address", url: "URL", emoji: "emoji", uuid: "UUID", uuidv4: "UUIDv4", uuidv6: "UUIDv6", nanoid: "nanoid", guid: "GUID", cuid: "cuid", cuid2: "cuid2", ulid: "ULID", xid: "XID", ksuid: "KSUID", datetime: "ISO datetime", date: "ISO date", time: "ISO time", duration: "ISO duration", ipv4: "IPv4 address", ipv6: "IPv6 address", cidrv4: "IPv4 range", cidrv6: "IPv6 range", base64: "base64-encoded string", base64url: "base64url-encoded string", json_string: "JSON string", e164: "E.164 number", jwt: "JWT", template_literal: "input" };
  return (Y) => {
    switch (Y.code) {
      case "invalid_type":
        return `Invalid input: expected ${Y.expected}, received ${YN(Y.input)}`;
      case "invalid_value":
        if (Y.values.length === 1) return `Invalid input: expected ${K8(Y.values[0])}`;
        return `Invalid option: expected one of ${B8(Y.values, "|")}`;
      case "too_big": {
        let J = Y.inclusive ? "<=" : "<", G = $(Y.origin);
        if (G) return `Too big: expected ${Y.origin ?? "value"} to have ${J}${Y.maximum.toString()} ${G.unit ?? "elements"}`;
        return `Too big: expected ${Y.origin ?? "value"} to be ${J}${Y.maximum.toString()}`;
      }
      case "too_small": {
        let J = Y.inclusive ? ">=" : ">", G = $(Y.origin);
        if (G) return `Too small: expected ${Y.origin} to have ${J}${Y.minimum.toString()} ${G.unit}`;
        return `Too small: expected ${Y.origin} to be ${J}${Y.minimum.toString()}`;
      }
      case "invalid_format": {
        let J = Y;
        if (J.format === "starts_with") return `Invalid string: must start with "${J.prefix}"`;
        if (J.format === "ends_with") return `Invalid string: must end with "${J.suffix}"`;
        if (J.format === "includes") return `Invalid string: must include "${J.includes}"`;
        if (J.format === "regex") return `Invalid string: must match pattern ${J.pattern}`;
        return `Invalid ${X[J.format] ?? Y.format}`;
      }
      case "not_multiple_of":
        return `Invalid number: must be a multiple of ${Y.divisor}`;
      case "unrecognized_keys":
        return `Unrecognized key${Y.keys.length > 1 ? "s" : ""}: ${B8(Y.keys, ", ")}`;
      case "invalid_key":
        return `Invalid key in ${Y.origin}`;
      case "invalid_union":
        return "Invalid input";
      case "invalid_element":
        return `Invalid value in ${Y.origin}`;
      default:
        return "Invalid input";
    }
  };
};
function wX() {
  return { localeError: JN() };
}
var O8 = class {
  constructor() {
    this._map = /* @__PURE__ */ new WeakMap(), this._idmap = /* @__PURE__ */ new Map();
  }
  add(Q, ...$) {
    let X = $[0];
    if (this._map.set(Q, X), X && typeof X === "object" && "id" in X) {
      if (this._idmap.has(X.id)) throw Error(`ID ${X.id} already exists in the registry`);
      this._idmap.set(X.id, Q);
    }
    return this;
  }
  remove(Q) {
    return this._map.delete(Q), this;
  }
  get(Q) {
    let $ = Q._zod.parent;
    if ($) {
      let X = { ...this.get($) ?? {} };
      return delete X.id, { ...X, ...this._map.get(Q) };
    }
    return this._map.get(Q);
  }
  has(Q) {
    return this._map.has(Q);
  }
};
function gG() {
  return new O8();
}
var r1 = gG();
function MX(Q, $) {
  return new Q({ type: "string", ...y($) });
}
function AX(Q, $) {
  return new Q({ type: "string", format: "email", check: "string_format", abort: false, ...y($) });
}
function D8(Q, $) {
  return new Q({ type: "string", format: "guid", check: "string_format", abort: false, ...y($) });
}
function jX(Q, $) {
  return new Q({ type: "string", format: "uuid", check: "string_format", abort: false, ...y($) });
}
function RX(Q, $) {
  return new Q({ type: "string", format: "uuid", check: "string_format", abort: false, version: "v4", ...y($) });
}
function IX(Q, $) {
  return new Q({ type: "string", format: "uuid", check: "string_format", abort: false, version: "v6", ...y($) });
}
function EX(Q, $) {
  return new Q({ type: "string", format: "uuid", check: "string_format", abort: false, version: "v7", ...y($) });
}
function PX(Q, $) {
  return new Q({ type: "string", format: "url", check: "string_format", abort: false, ...y($) });
}
function bX(Q, $) {
  return new Q({ type: "string", format: "emoji", check: "string_format", abort: false, ...y($) });
}
function ZX(Q, $) {
  return new Q({ type: "string", format: "nanoid", check: "string_format", abort: false, ...y($) });
}
function CX(Q, $) {
  return new Q({ type: "string", format: "cuid", check: "string_format", abort: false, ...y($) });
}
function SX(Q, $) {
  return new Q({ type: "string", format: "cuid2", check: "string_format", abort: false, ...y($) });
}
function _X(Q, $) {
  return new Q({ type: "string", format: "ulid", check: "string_format", abort: false, ...y($) });
}
function kX(Q, $) {
  return new Q({ type: "string", format: "xid", check: "string_format", abort: false, ...y($) });
}
function vX(Q, $) {
  return new Q({ type: "string", format: "ksuid", check: "string_format", abort: false, ...y($) });
}
function TX(Q, $) {
  return new Q({ type: "string", format: "ipv4", check: "string_format", abort: false, ...y($) });
}
function xX(Q, $) {
  return new Q({ type: "string", format: "ipv6", check: "string_format", abort: false, ...y($) });
}
function yX(Q, $) {
  return new Q({ type: "string", format: "cidrv4", check: "string_format", abort: false, ...y($) });
}
function gX(Q, $) {
  return new Q({ type: "string", format: "cidrv6", check: "string_format", abort: false, ...y($) });
}
function hX(Q, $) {
  return new Q({ type: "string", format: "base64", check: "string_format", abort: false, ...y($) });
}
function fX(Q, $) {
  return new Q({ type: "string", format: "base64url", check: "string_format", abort: false, ...y($) });
}
function uX(Q, $) {
  return new Q({ type: "string", format: "e164", check: "string_format", abort: false, ...y($) });
}
function mX(Q, $) {
  return new Q({ type: "string", format: "jwt", check: "string_format", abort: false, ...y($) });
}
function hG(Q, $) {
  return new Q({ type: "string", format: "datetime", check: "string_format", offset: false, local: false, precision: null, ...y($) });
}
function fG(Q, $) {
  return new Q({ type: "string", format: "date", check: "string_format", ...y($) });
}
function uG(Q, $) {
  return new Q({ type: "string", format: "time", check: "string_format", precision: null, ...y($) });
}
function mG(Q, $) {
  return new Q({ type: "string", format: "duration", check: "string_format", ...y($) });
}
function lX(Q, $) {
  return new Q({ type: "number", checks: [], ...y($) });
}
function cX(Q, $) {
  return new Q({ type: "number", check: "number_format", abort: false, format: "safeint", ...y($) });
}
function pX(Q, $) {
  return new Q({ type: "boolean", ...y($) });
}
function dX(Q, $) {
  return new Q({ type: "null", ...y($) });
}
function iX(Q) {
  return new Q({ type: "unknown" });
}
function nX(Q, $) {
  return new Q({ type: "never", ...y($) });
}
function w8(Q, $) {
  return new C$({ check: "less_than", ...y($), value: Q, inclusive: false });
}
function J9(Q, $) {
  return new C$({ check: "less_than", ...y($), value: Q, inclusive: true });
}
function M8(Q, $) {
  return new S$({ check: "greater_than", ...y($), value: Q, inclusive: false });
}
function G9(Q, $) {
  return new S$({ check: "greater_than", ...y($), value: Q, inclusive: true });
}
function A8(Q, $) {
  return new zG({ check: "multiple_of", ...y($), value: Q });
}
function j8(Q, $) {
  return new VG({ check: "max_length", ...y($), maximum: Q });
}
function q4(Q, $) {
  return new qG({ check: "min_length", ...y($), minimum: Q });
}
function R8(Q, $) {
  return new UG({ check: "length_equals", ...y($), length: Q });
}
function oX(Q, $) {
  return new LG({ check: "string_format", format: "regex", ...y($), pattern: Q });
}
function rX(Q) {
  return new FG({ check: "string_format", format: "lowercase", ...y(Q) });
}
function tX(Q) {
  return new NG({ check: "string_format", format: "uppercase", ...y(Q) });
}
function aX(Q, $) {
  return new OG({ check: "string_format", format: "includes", ...y($), includes: Q });
}
function sX(Q, $) {
  return new DG({ check: "string_format", format: "starts_with", ...y($), prefix: Q });
}
function eX(Q, $) {
  return new wG({ check: "string_format", format: "ends_with", ...y($), suffix: Q });
}
function I6(Q) {
  return new MG({ check: "overwrite", tx: Q });
}
function QY(Q) {
  return I6(($) => $.normalize(Q));
}
function $Y() {
  return I6((Q) => Q.trim());
}
function XY() {
  return I6((Q) => Q.toLowerCase());
}
function YY() {
  return I6((Q) => Q.toUpperCase());
}
function lG(Q, $, X) {
  return new Q({ type: "array", element: $, ...y(X) });
}
function JY(Q, $, X) {
  let Y = y(X);
  return Y.abort ?? (Y.abort = true), new Q({ type: "custom", check: "custom", fn: $, ...Y });
}
function GY(Q, $, X) {
  return new Q({ type: "custom", check: "custom", fn: $, ...y(X) });
}
var mN = O("ZodMiniType", (Q, $) => {
  if (!Q._zod) throw Error("Uninitialized schema in ZodMiniType.");
  e.init(Q, $), Q.def = $, Q.parse = (X, Y) => j$(Q, X, Y, { callee: Q.parse }), Q.safeParse = (X, Y) => j6(Q, X, Y), Q.parseAsync = async (X, Y) => I$(Q, X, Y, { callee: Q.parseAsync }), Q.safeParseAsync = async (X, Y) => R6(Q, X, Y), Q.check = (...X) => {
    return Q.clone({ ...$, checks: [...$.checks ?? [], ...X.map((Y) => typeof Y === "function" ? { _zod: { check: Y, def: { check: "custom" }, onattach: [] } } : Y)] });
  }, Q.clone = (X, Y) => c0(Q, X, Y), Q.brand = () => Q, Q.register = (X, Y) => {
    return X.add(Q, Y), Q;
  };
});
var lN = O("ZodMiniObject", (Q, $) => {
  F8.init(Q, $), mN.init(Q, $), i.defineLazy(Q, "shape", () => $.shape);
});
var W9 = {};
uQ(W9, { time: () => VY, duration: () => qY, datetime: () => zY, date: () => KY, ZodISOTime: () => oG, ZodISODuration: () => rG, ZodISODateTime: () => iG, ZodISODate: () => nG });
var iG = O("ZodISODateTime", (Q, $) => {
  kG.init(Q, $), B0.init(Q, $);
});
function zY(Q) {
  return hG(iG, Q);
}
var nG = O("ZodISODate", (Q, $) => {
  vG.init(Q, $), B0.init(Q, $);
});
function KY(Q) {
  return fG(nG, Q);
}
var oG = O("ZodISOTime", (Q, $) => {
  TG.init(Q, $), B0.init(Q, $);
});
function VY(Q) {
  return uG(oG, Q);
}
var rG = O("ZodISODuration", (Q, $) => {
  xG.init(Q, $), B0.init(Q, $);
});
function qY(Q) {
  return mG(rG, Q);
}
var tG = (Q, $) => {
  V8.init(Q, $), Q.name = "ZodError", Object.defineProperties(Q, { format: { value: (X) => M$(Q, X) }, flatten: { value: (X) => w$(Q, X) }, addIssue: { value: (X) => Q.issues.push(X) }, addIssues: { value: (X) => Q.issues.push(...X) }, isEmpty: { get() {
    return Q.issues.length === 0;
  } } });
};
var Fk = O("ZodError", tG);
var H9 = O("ZodError", tG, { Parent: Error });
var aG = A$(H9);
var sG = R$(H9);
var eG = E$(H9);
var QW = P$(H9);
var F0 = O("ZodType", (Q, $) => {
  return e.init(Q, $), Q.def = $, Object.defineProperty(Q, "_def", { value: $ }), Q.check = (...X) => {
    return Q.clone({ ...$, checks: [...$.checks ?? [], ...X.map((Y) => typeof Y === "function" ? { _zod: { check: Y, def: { check: "custom" }, onattach: [] } } : Y)] });
  }, Q.clone = (X, Y) => c0(Q, X, Y), Q.brand = () => Q, Q.register = (X, Y) => {
    return X.add(Q, Y), Q;
  }, Q.parse = (X, Y) => aG(Q, X, Y, { callee: Q.parse }), Q.safeParse = (X, Y) => eG(Q, X, Y), Q.parseAsync = async (X, Y) => sG(Q, X, Y, { callee: Q.parseAsync }), Q.safeParseAsync = async (X, Y) => QW(Q, X, Y), Q.spa = Q.safeParseAsync, Q.refine = (X, Y) => Q.check(fO(X, Y)), Q.superRefine = (X) => Q.check(uO(X)), Q.overwrite = (X) => Q.check(I6(X)), Q.optional = () => L0(Q), Q.nullable = () => YW(Q), Q.nullish = () => L0(YW(Q)), Q.nonoptional = (X) => kO(Q, X), Q.array = () => n(Q), Q.or = (X) => G0([Q, X]), Q.and = (X) => Z8(Q, X), Q.transform = (X) => LY(Q, BW(X)), Q.default = (X) => CO(Q, X), Q.prefault = (X) => _O(Q, X), Q.catch = (X) => TO(Q, X), Q.pipe = (X) => LY(Q, X), Q.readonly = () => gO(Q), Q.describe = (X) => {
    let Y = Q.clone();
    return r1.add(Y, { description: X }), Y;
  }, Object.defineProperty(Q, "description", { get() {
    return r1.get(Q)?.description;
  }, configurable: true }), Q.meta = (...X) => {
    if (X.length === 0) return r1.get(Q);
    let Y = Q.clone();
    return r1.add(Y, X[0]), Y;
  }, Q.isOptional = () => Q.safeParse(void 0).success, Q.isNullable = () => Q.safeParse(null).success, Q;
});
var JW = O("_ZodString", (Q, $) => {
  Y9.init(Q, $), F0.init(Q, $);
  let X = Q._zod.bag;
  Q.format = X.format ?? null, Q.minLength = X.minimum ?? null, Q.maxLength = X.maximum ?? null, Q.regex = (...Y) => Q.check(oX(...Y)), Q.includes = (...Y) => Q.check(aX(...Y)), Q.startsWith = (...Y) => Q.check(sX(...Y)), Q.endsWith = (...Y) => Q.check(eX(...Y)), Q.min = (...Y) => Q.check(q4(...Y)), Q.max = (...Y) => Q.check(j8(...Y)), Q.length = (...Y) => Q.check(R8(...Y)), Q.nonempty = (...Y) => Q.check(q4(1, ...Y)), Q.lowercase = (Y) => Q.check(rX(Y)), Q.uppercase = (Y) => Q.check(tX(Y)), Q.trim = () => Q.check($Y()), Q.normalize = (...Y) => Q.check(QY(...Y)), Q.toLowerCase = () => Q.check(XY()), Q.toUpperCase = () => Q.check(YY());
});
var aN = O("ZodString", (Q, $) => {
  Y9.init(Q, $), JW.init(Q, $), Q.email = (X) => Q.check(AX(sN, X)), Q.url = (X) => Q.check(PX(eN, X)), Q.jwt = (X) => Q.check(mX(LO, X)), Q.emoji = (X) => Q.check(bX(QO, X)), Q.guid = (X) => Q.check(D8($W, X)), Q.uuid = (X) => Q.check(jX(b8, X)), Q.uuidv4 = (X) => Q.check(RX(b8, X)), Q.uuidv6 = (X) => Q.check(IX(b8, X)), Q.uuidv7 = (X) => Q.check(EX(b8, X)), Q.nanoid = (X) => Q.check(ZX($O, X)), Q.guid = (X) => Q.check(D8($W, X)), Q.cuid = (X) => Q.check(CX(XO, X)), Q.cuid2 = (X) => Q.check(SX(YO, X)), Q.ulid = (X) => Q.check(_X(JO, X)), Q.base64 = (X) => Q.check(hX(VO, X)), Q.base64url = (X) => Q.check(fX(qO, X)), Q.xid = (X) => Q.check(kX(GO, X)), Q.ksuid = (X) => Q.check(vX(WO, X)), Q.ipv4 = (X) => Q.check(TX(HO, X)), Q.ipv6 = (X) => Q.check(xX(BO, X)), Q.cidrv4 = (X) => Q.check(yX(zO, X)), Q.cidrv6 = (X) => Q.check(gX(KO, X)), Q.e164 = (X) => Q.check(uX(UO, X)), Q.datetime = (X) => Q.check(zY(X)), Q.date = (X) => Q.check(KY(X)), Q.time = (X) => Q.check(VY(X)), Q.duration = (X) => Q.check(qY(X));
});
function N(Q) {
  return MX(aN, Q);
}
var B0 = O("ZodStringFormat", (Q, $) => {
  J0.init(Q, $), JW.init(Q, $);
});
var sN = O("ZodEmail", (Q, $) => {
  x$.init(Q, $), B0.init(Q, $);
});
var $W = O("ZodGUID", (Q, $) => {
  v$.init(Q, $), B0.init(Q, $);
});
var b8 = O("ZodUUID", (Q, $) => {
  T$.init(Q, $), B0.init(Q, $);
});
var eN = O("ZodURL", (Q, $) => {
  y$.init(Q, $), B0.init(Q, $);
});
var QO = O("ZodEmoji", (Q, $) => {
  g$.init(Q, $), B0.init(Q, $);
});
var $O = O("ZodNanoID", (Q, $) => {
  h$.init(Q, $), B0.init(Q, $);
});
var XO = O("ZodCUID", (Q, $) => {
  f$.init(Q, $), B0.init(Q, $);
});
var YO = O("ZodCUID2", (Q, $) => {
  u$.init(Q, $), B0.init(Q, $);
});
var JO = O("ZodULID", (Q, $) => {
  m$.init(Q, $), B0.init(Q, $);
});
var GO = O("ZodXID", (Q, $) => {
  l$.init(Q, $), B0.init(Q, $);
});
var WO = O("ZodKSUID", (Q, $) => {
  c$.init(Q, $), B0.init(Q, $);
});
var HO = O("ZodIPv4", (Q, $) => {
  p$.init(Q, $), B0.init(Q, $);
});
var BO = O("ZodIPv6", (Q, $) => {
  d$.init(Q, $), B0.init(Q, $);
});
var zO = O("ZodCIDRv4", (Q, $) => {
  i$.init(Q, $), B0.init(Q, $);
});
var KO = O("ZodCIDRv6", (Q, $) => {
  n$.init(Q, $), B0.init(Q, $);
});
var VO = O("ZodBase64", (Q, $) => {
  o$.init(Q, $), B0.init(Q, $);
});
var qO = O("ZodBase64URL", (Q, $) => {
  r$.init(Q, $), B0.init(Q, $);
});
var UO = O("ZodE164", (Q, $) => {
  t$.init(Q, $), B0.init(Q, $);
});
var LO = O("ZodJWT", (Q, $) => {
  a$.init(Q, $), B0.init(Q, $);
});
var GW = O("ZodNumber", (Q, $) => {
  L8.init(Q, $), F0.init(Q, $), Q.gt = (Y, J) => Q.check(M8(Y, J)), Q.gte = (Y, J) => Q.check(G9(Y, J)), Q.min = (Y, J) => Q.check(G9(Y, J)), Q.lt = (Y, J) => Q.check(w8(Y, J)), Q.lte = (Y, J) => Q.check(J9(Y, J)), Q.max = (Y, J) => Q.check(J9(Y, J)), Q.int = (Y) => Q.check(XW(Y)), Q.safe = (Y) => Q.check(XW(Y)), Q.positive = (Y) => Q.check(M8(0, Y)), Q.nonnegative = (Y) => Q.check(G9(0, Y)), Q.negative = (Y) => Q.check(w8(0, Y)), Q.nonpositive = (Y) => Q.check(J9(0, Y)), Q.multipleOf = (Y, J) => Q.check(A8(Y, J)), Q.step = (Y, J) => Q.check(A8(Y, J)), Q.finite = () => Q;
  let X = Q._zod.bag;
  Q.minValue = Math.max(X.minimum ?? Number.NEGATIVE_INFINITY, X.exclusiveMinimum ?? Number.NEGATIVE_INFINITY) ?? null, Q.maxValue = Math.min(X.maximum ?? Number.POSITIVE_INFINITY, X.exclusiveMaximum ?? Number.POSITIVE_INFINITY) ?? null, Q.isInt = (X.format ?? "").includes("int") || Number.isSafeInteger(X.multipleOf ?? 0.5), Q.isFinite = true, Q.format = X.format ?? null;
});
function s(Q) {
  return lX(GW, Q);
}
var FO = O("ZodNumberFormat", (Q, $) => {
  s$.init(Q, $), GW.init(Q, $);
});
function XW(Q) {
  return cX(FO, Q);
}
var NO = O("ZodBoolean", (Q, $) => {
  e$.init(Q, $), F0.init(Q, $);
});
function M0(Q) {
  return pX(NO, Q);
}
var OO = O("ZodNull", (Q, $) => {
  QX.init(Q, $), F0.init(Q, $);
});
function FY(Q) {
  return dX(OO, Q);
}
var DO = O("ZodUnknown", (Q, $) => {
  $X.init(Q, $), F0.init(Q, $);
});
function z0() {
  return iX(DO);
}
var wO = O("ZodNever", (Q, $) => {
  XX.init(Q, $), F0.init(Q, $);
});
function MO(Q) {
  return nX(wO, Q);
}
var AO = O("ZodArray", (Q, $) => {
  YX.init(Q, $), F0.init(Q, $), Q.element = $.element, Q.min = (X, Y) => Q.check(q4(X, Y)), Q.nonempty = (X) => Q.check(q4(1, X)), Q.max = (X, Y) => Q.check(j8(X, Y)), Q.length = (X, Y) => Q.check(R8(X, Y)), Q.unwrap = () => Q.element;
});
function n(Q, $) {
  return lG(AO, Q, $);
}
var WW = O("ZodObject", (Q, $) => {
  F8.init(Q, $), F0.init(Q, $), i.defineLazy(Q, "shape", () => $.shape), Q.keyof = () => y0(Object.keys(Q._zod.def.shape)), Q.catchall = (X) => Q.clone({ ...Q._zod.def, catchall: X }), Q.passthrough = () => Q.clone({ ...Q._zod.def, catchall: z0() }), Q.loose = () => Q.clone({ ...Q._zod.def, catchall: z0() }), Q.strict = () => Q.clone({ ...Q._zod.def, catchall: MO() }), Q.strip = () => Q.clone({ ...Q._zod.def, catchall: void 0 }), Q.extend = (X) => {
    return i.extend(Q, X);
  }, Q.merge = (X) => i.merge(Q, X), Q.pick = (X) => i.pick(Q, X), Q.omit = (X) => i.omit(Q, X), Q.partial = (...X) => i.partial(zW, Q, X[0]), Q.required = (...X) => i.required(KW, Q, X[0]);
});
function b(Q, $) {
  let X = { type: "object", get shape() {
    return i.assignProp(this, "shape", { ...Q }), this.shape;
  }, ...i.normalizeParams($) };
  return new WW(X);
}
function S0(Q, $) {
  return new WW({ type: "object", get shape() {
    return i.assignProp(this, "shape", { ...Q }), this.shape;
  }, catchall: z0(), ...i.normalizeParams($) });
}
var HW = O("ZodUnion", (Q, $) => {
  N8.init(Q, $), F0.init(Q, $), Q.options = $.options;
});
function G0(Q, $) {
  return new HW({ type: "union", options: Q, ...i.normalizeParams($) });
}
var jO = O("ZodDiscriminatedUnion", (Q, $) => {
  HW.init(Q, $), JX.init(Q, $);
});
function NY(Q, $, X) {
  return new jO({ type: "union", options: $, discriminator: Q, ...i.normalizeParams(X) });
}
var RO = O("ZodIntersection", (Q, $) => {
  GX.init(Q, $), F0.init(Q, $);
});
function Z8(Q, $) {
  return new RO({ type: "intersection", left: Q, right: $ });
}
var IO = O("ZodRecord", (Q, $) => {
  WX.init(Q, $), F0.init(Q, $), Q.keyType = $.keyType, Q.valueType = $.valueType;
});
function K0(Q, $, X) {
  return new IO({ type: "record", keyType: Q, valueType: $, ...i.normalizeParams(X) });
}
var UY = O("ZodEnum", (Q, $) => {
  HX.init(Q, $), F0.init(Q, $), Q.enum = $.entries, Q.options = Object.values($.entries);
  let X = new Set(Object.keys($.entries));
  Q.extract = (Y, J) => {
    let G = {};
    for (let W of Y) if (X.has(W)) G[W] = $.entries[W];
    else throw Error(`Key ${W} not found in enum`);
    return new UY({ ...$, checks: [], ...i.normalizeParams(J), entries: G });
  }, Q.exclude = (Y, J) => {
    let G = { ...$.entries };
    for (let W of Y) if (X.has(W)) delete G[W];
    else throw Error(`Key ${W} not found in enum`);
    return new UY({ ...$, checks: [], ...i.normalizeParams(J), entries: G });
  };
});
function y0(Q, $) {
  let X = Array.isArray(Q) ? Object.fromEntries(Q.map((Y) => [Y, Y])) : Q;
  return new UY({ type: "enum", entries: X, ...i.normalizeParams($) });
}
var EO = O("ZodLiteral", (Q, $) => {
  BX.init(Q, $), F0.init(Q, $), Q.values = new Set($.values), Object.defineProperty(Q, "value", { get() {
    if ($.values.length > 1) throw Error("This schema contains multiple valid literal values. Use `.values` instead.");
    return $.values[0];
  } });
});
function _(Q, $) {
  return new EO({ type: "literal", values: Array.isArray(Q) ? Q : [Q], ...i.normalizeParams($) });
}
var PO = O("ZodTransform", (Q, $) => {
  zX.init(Q, $), F0.init(Q, $), Q._zod.parse = (X, Y) => {
    X.addIssue = (G) => {
      if (typeof G === "string") X.issues.push(i.issue(G, X.value, $));
      else {
        let W = G;
        if (W.fatal) W.continue = false;
        W.code ?? (W.code = "custom"), W.input ?? (W.input = X.value), W.inst ?? (W.inst = Q), W.continue ?? (W.continue = true), X.issues.push(i.issue(W));
      }
    };
    let J = $.transform(X.value, X);
    if (J instanceof Promise) return J.then((G) => {
      return X.value = G, X;
    });
    return X.value = J, X;
  };
});
function BW(Q) {
  return new PO({ type: "transform", transform: Q });
}
var zW = O("ZodOptional", (Q, $) => {
  KX.init(Q, $), F0.init(Q, $), Q.unwrap = () => Q._zod.def.innerType;
});
function L0(Q) {
  return new zW({ type: "optional", innerType: Q });
}
var bO = O("ZodNullable", (Q, $) => {
  VX.init(Q, $), F0.init(Q, $), Q.unwrap = () => Q._zod.def.innerType;
});
function YW(Q) {
  return new bO({ type: "nullable", innerType: Q });
}
var ZO = O("ZodDefault", (Q, $) => {
  qX.init(Q, $), F0.init(Q, $), Q.unwrap = () => Q._zod.def.innerType, Q.removeDefault = Q.unwrap;
});
function CO(Q, $) {
  return new ZO({ type: "default", innerType: Q, get defaultValue() {
    return typeof $ === "function" ? $() : $;
  } });
}
var SO = O("ZodPrefault", (Q, $) => {
  UX.init(Q, $), F0.init(Q, $), Q.unwrap = () => Q._zod.def.innerType;
});
function _O(Q, $) {
  return new SO({ type: "prefault", innerType: Q, get defaultValue() {
    return typeof $ === "function" ? $() : $;
  } });
}
var KW = O("ZodNonOptional", (Q, $) => {
  LX.init(Q, $), F0.init(Q, $), Q.unwrap = () => Q._zod.def.innerType;
});
function kO(Q, $) {
  return new KW({ type: "nonoptional", innerType: Q, ...i.normalizeParams($) });
}
var vO = O("ZodCatch", (Q, $) => {
  FX.init(Q, $), F0.init(Q, $), Q.unwrap = () => Q._zod.def.innerType, Q.removeCatch = Q.unwrap;
});
function TO(Q, $) {
  return new vO({ type: "catch", innerType: Q, catchValue: typeof $ === "function" ? $ : () => $ });
}
var xO = O("ZodPipe", (Q, $) => {
  NX.init(Q, $), F0.init(Q, $), Q.in = $.in, Q.out = $.out;
});
function LY(Q, $) {
  return new xO({ type: "pipe", in: Q, out: $ });
}
var yO = O("ZodReadonly", (Q, $) => {
  OX.init(Q, $), F0.init(Q, $);
});
function gO(Q) {
  return new yO({ type: "readonly", innerType: Q });
}
var VW = O("ZodCustom", (Q, $) => {
  DX.init(Q, $), F0.init(Q, $);
});
function hO(Q, $) {
  let X = new j0({ check: "custom", ...i.normalizeParams($) });
  return X._zod.check = Q, X;
}
function qW(Q, $) {
  return JY(VW, Q ?? (() => true), $);
}
function fO(Q, $ = {}) {
  return GY(VW, Q, $);
}
function uO(Q, $) {
  let X = hO((Y) => {
    return Y.addIssue = (J) => {
      if (typeof J === "string") Y.issues.push(i.issue(J, Y.value, X._zod.def));
      else {
        let G = J;
        if (G.fatal) G.continue = false;
        G.code ?? (G.code = "custom"), G.input ?? (G.input = Y.value), G.inst ?? (G.inst = X), G.continue ?? (G.continue = !X._zod.def.abort), Y.issues.push(i.issue(G));
      }
    }, Q(Y.value, Y);
  }, $);
  return X;
}
function OY(Q, $) {
  return LY(BW(Q), $);
}
l0(wX());
var s1 = "io.modelcontextprotocol/related-task";
var S8 = "2.0";
var R0 = qW((Q) => Q !== null && (typeof Q === "object" || typeof Q === "function"));
var LW = G0([N(), s().int()]);
var FW = N();
var Ck = S0({ ttl: G0([s(), FY()]).optional(), pollInterval: s().optional() });
var mO = b({ ttl: s().optional() });
var lO = b({ taskId: N() });
var wY = S0({ progressToken: LW.optional(), [s1]: lO.optional() });
var p0 = b({ _meta: wY.optional() });
var B9 = p0.extend({ task: mO.optional() });
var I0 = b({ method: N(), params: p0.loose().optional() });
var a0 = b({ _meta: wY.optional() });
var s0 = b({ method: N(), params: a0.loose().optional() });
var E0 = S0({ _meta: wY.optional() });
var _8 = G0([N(), s().int()]);
var OW = b({ jsonrpc: _(S8), id: _8, ...I0.shape }).strict();
var DW = b({ jsonrpc: _(S8), ...s0.shape }).strict();
var AY = b({ jsonrpc: _(S8), id: _8, result: E0 }).strict();
var T;
(function(Q) {
  Q[Q.ConnectionClosed = -32e3] = "ConnectionClosed", Q[Q.RequestTimeout = -32001] = "RequestTimeout", Q[Q.ParseError = -32700] = "ParseError", Q[Q.InvalidRequest = -32600] = "InvalidRequest", Q[Q.MethodNotFound = -32601] = "MethodNotFound", Q[Q.InvalidParams = -32602] = "InvalidParams", Q[Q.InternalError = -32603] = "InternalError", Q[Q.UrlElicitationRequired = -32042] = "UrlElicitationRequired";
})(T || (T = {}));
var jY = b({ jsonrpc: _(S8), id: _8.optional(), error: b({ code: s().int(), message: N(), data: z0().optional() }) }).strict();
var Sk = G0([OW, DW, AY, jY]);
var _k = G0([AY, jY]);
var k8 = E0.strict();
var cO = a0.extend({ requestId: _8.optional(), reason: N().optional() });
var v8 = s0.extend({ method: _("notifications/cancelled"), params: cO });
var pO = b({ src: N(), mimeType: N().optional(), sizes: n(N()).optional(), theme: y0(["light", "dark"]).optional() });
var K9 = b({ icons: n(pO).optional() });
var L4 = b({ name: N(), title: N().optional() });
var AW = L4.extend({ ...L4.shape, ...K9.shape, version: N(), websiteUrl: N().optional(), description: N().optional() });
var dO = Z8(b({ applyDefaults: M0().optional() }), K0(N(), z0()));
var iO = OY((Q) => {
  if (Q && typeof Q === "object" && !Array.isArray(Q)) {
    if (Object.keys(Q).length === 0) return { form: {} };
  }
  return Q;
}, Z8(b({ form: dO.optional(), url: R0.optional() }), K0(N(), z0()).optional()));
var nO = S0({ list: R0.optional(), cancel: R0.optional(), requests: S0({ sampling: S0({ createMessage: R0.optional() }).optional(), elicitation: S0({ create: R0.optional() }).optional() }).optional() });
var oO = S0({ list: R0.optional(), cancel: R0.optional(), requests: S0({ tools: S0({ call: R0.optional() }).optional() }).optional() });
var rO = b({ experimental: K0(N(), R0).optional(), sampling: b({ context: R0.optional(), tools: R0.optional() }).optional(), elicitation: iO.optional(), roots: b({ listChanged: M0().optional() }).optional(), tasks: nO.optional() });
var tO = p0.extend({ protocolVersion: N(), capabilities: rO, clientInfo: AW });
var RY = I0.extend({ method: _("initialize"), params: tO });
var aO = b({ experimental: K0(N(), R0).optional(), logging: R0.optional(), completions: R0.optional(), prompts: b({ listChanged: M0().optional() }).optional(), resources: b({ subscribe: M0().optional(), listChanged: M0().optional() }).optional(), tools: b({ listChanged: M0().optional() }).optional(), tasks: oO.optional() });
var sO = E0.extend({ protocolVersion: N(), capabilities: aO, serverInfo: AW, instructions: N().optional() });
var IY = s0.extend({ method: _("notifications/initialized"), params: a0.optional() });
var T8 = I0.extend({ method: _("ping"), params: p0.optional() });
var eO = b({ progress: s(), total: L0(s()), message: L0(N()) });
var QD = b({ ...a0.shape, ...eO.shape, progressToken: LW });
var x8 = s0.extend({ method: _("notifications/progress"), params: QD });
var $D = p0.extend({ cursor: FW.optional() });
var V9 = I0.extend({ params: $D.optional() });
var q9 = E0.extend({ nextCursor: FW.optional() });
var XD = y0(["working", "input_required", "completed", "failed", "cancelled"]);
var U9 = b({ taskId: N(), status: XD, ttl: G0([s(), FY()]), createdAt: N(), lastUpdatedAt: N(), pollInterval: L0(s()), statusMessage: L0(N()) });
var F4 = E0.extend({ task: U9 });
var YD = a0.merge(U9);
var L9 = s0.extend({ method: _("notifications/tasks/status"), params: YD });
var y8 = I0.extend({ method: _("tasks/get"), params: p0.extend({ taskId: N() }) });
var g8 = E0.merge(U9);
var h8 = I0.extend({ method: _("tasks/result"), params: p0.extend({ taskId: N() }) });
var kk = E0.loose();
var f8 = V9.extend({ method: _("tasks/list") });
var u8 = q9.extend({ tasks: n(U9) });
var m8 = I0.extend({ method: _("tasks/cancel"), params: p0.extend({ taskId: N() }) });
var jW = E0.merge(U9);
var RW = b({ uri: N(), mimeType: L0(N()), _meta: K0(N(), z0()).optional() });
var IW = RW.extend({ text: N() });
var EY = N().refine((Q) => {
  try {
    return atob(Q), true;
  } catch {
    return false;
  }
}, { message: "Invalid Base64 string" });
var EW = RW.extend({ blob: EY });
var F9 = y0(["user", "assistant"]);
var N4 = b({ audience: n(F9).optional(), priority: s().min(0).max(1).optional(), lastModified: W9.datetime({ offset: true }).optional() });
var PW = b({ ...L4.shape, ...K9.shape, uri: N(), description: L0(N()), mimeType: L0(N()), annotations: N4.optional(), _meta: L0(S0({})) });
var JD = b({ ...L4.shape, ...K9.shape, uriTemplate: N(), description: L0(N()), mimeType: L0(N()), annotations: N4.optional(), _meta: L0(S0({})) });
var l8 = V9.extend({ method: _("resources/list") });
var GD = q9.extend({ resources: n(PW) });
var c8 = V9.extend({ method: _("resources/templates/list") });
var WD = q9.extend({ resourceTemplates: n(JD) });
var PY = p0.extend({ uri: N() });
var HD = PY;
var p8 = I0.extend({ method: _("resources/read"), params: HD });
var BD = E0.extend({ contents: n(G0([IW, EW])) });
var zD = s0.extend({ method: _("notifications/resources/list_changed"), params: a0.optional() });
var KD = PY;
var VD = I0.extend({ method: _("resources/subscribe"), params: KD });
var qD = PY;
var UD = I0.extend({ method: _("resources/unsubscribe"), params: qD });
var LD = a0.extend({ uri: N() });
var FD = s0.extend({ method: _("notifications/resources/updated"), params: LD });
var ND = b({ name: N(), description: L0(N()), required: L0(M0()) });
var OD = b({ ...L4.shape, ...K9.shape, description: L0(N()), arguments: L0(n(ND)), _meta: L0(S0({})) });
var d8 = V9.extend({ method: _("prompts/list") });
var DD = q9.extend({ prompts: n(OD) });
var wD = p0.extend({ name: N(), arguments: K0(N(), N()).optional() });
var i8 = I0.extend({ method: _("prompts/get"), params: wD });
var bY = b({ type: _("text"), text: N(), annotations: N4.optional(), _meta: K0(N(), z0()).optional() });
var ZY = b({ type: _("image"), data: EY, mimeType: N(), annotations: N4.optional(), _meta: K0(N(), z0()).optional() });
var CY = b({ type: _("audio"), data: EY, mimeType: N(), annotations: N4.optional(), _meta: K0(N(), z0()).optional() });
var MD = b({ type: _("tool_use"), name: N(), id: N(), input: K0(N(), z0()), _meta: K0(N(), z0()).optional() });
var AD = b({ type: _("resource"), resource: G0([IW, EW]), annotations: N4.optional(), _meta: K0(N(), z0()).optional() });
var jD = PW.extend({ type: _("resource_link") });
var SY = G0([bY, ZY, CY, jD, AD]);
var RD = b({ role: F9, content: SY });
var ID = E0.extend({ description: N().optional(), messages: n(RD) });
var ED = s0.extend({ method: _("notifications/prompts/list_changed"), params: a0.optional() });
var PD = b({ title: N().optional(), readOnlyHint: M0().optional(), destructiveHint: M0().optional(), idempotentHint: M0().optional(), openWorldHint: M0().optional() });
var bD = b({ taskSupport: y0(["required", "optional", "forbidden"]).optional() });
var bW = b({ ...L4.shape, ...K9.shape, description: N().optional(), inputSchema: b({ type: _("object"), properties: K0(N(), R0).optional(), required: n(N()).optional() }).catchall(z0()), outputSchema: b({ type: _("object"), properties: K0(N(), R0).optional(), required: n(N()).optional() }).catchall(z0()).optional(), annotations: PD.optional(), execution: bD.optional(), _meta: K0(N(), z0()).optional() });
var n8 = V9.extend({ method: _("tools/list") });
var ZD = q9.extend({ tools: n(bW) });
var o8 = E0.extend({ content: n(SY).default([]), structuredContent: K0(N(), z0()).optional(), isError: M0().optional() });
var vk = o8.or(E0.extend({ toolResult: z0() }));
var CD = B9.extend({ name: N(), arguments: K0(N(), z0()).optional() });
var O4 = I0.extend({ method: _("tools/call"), params: CD });
var SD = s0.extend({ method: _("notifications/tools/list_changed"), params: a0.optional() });
var Tk = b({ autoRefresh: M0().default(true), debounceMs: s().int().nonnegative().default(300) });
var N9 = y0(["debug", "info", "notice", "warning", "error", "critical", "alert", "emergency"]);
var _D = p0.extend({ level: N9 });
var _Y = I0.extend({ method: _("logging/setLevel"), params: _D });
var kD = a0.extend({ level: N9, logger: N().optional(), data: z0() });
var vD = s0.extend({ method: _("notifications/message"), params: kD });
var TD = b({ name: N().optional() });
var xD = b({ hints: n(TD).optional(), costPriority: s().min(0).max(1).optional(), speedPriority: s().min(0).max(1).optional(), intelligencePriority: s().min(0).max(1).optional() });
var yD = b({ mode: y0(["auto", "required", "none"]).optional() });
var gD = b({ type: _("tool_result"), toolUseId: N().describe("The unique identifier for the corresponding tool call."), content: n(SY).default([]), structuredContent: b({}).loose().optional(), isError: M0().optional(), _meta: K0(N(), z0()).optional() });
var hD = NY("type", [bY, ZY, CY]);
var C8 = NY("type", [bY, ZY, CY, MD, gD]);
var fD = b({ role: F9, content: G0([C8, n(C8)]), _meta: K0(N(), z0()).optional() });
var uD = B9.extend({ messages: n(fD), modelPreferences: xD.optional(), systemPrompt: N().optional(), includeContext: y0(["none", "thisServer", "allServers"]).optional(), temperature: s().optional(), maxTokens: s().int(), stopSequences: n(N()).optional(), metadata: R0.optional(), tools: n(bW).optional(), toolChoice: yD.optional() });
var mD = I0.extend({ method: _("sampling/createMessage"), params: uD });
var O9 = E0.extend({ model: N(), stopReason: L0(y0(["endTurn", "stopSequence", "maxTokens"]).or(N())), role: F9, content: hD });
var kY = E0.extend({ model: N(), stopReason: L0(y0(["endTurn", "stopSequence", "maxTokens", "toolUse"]).or(N())), role: F9, content: G0([C8, n(C8)]) });
var lD = b({ type: _("boolean"), title: N().optional(), description: N().optional(), default: M0().optional() });
var cD = b({ type: _("string"), title: N().optional(), description: N().optional(), minLength: s().optional(), maxLength: s().optional(), format: y0(["email", "uri", "date", "date-time"]).optional(), default: N().optional() });
var pD = b({ type: y0(["number", "integer"]), title: N().optional(), description: N().optional(), minimum: s().optional(), maximum: s().optional(), default: s().optional() });
var dD = b({ type: _("string"), title: N().optional(), description: N().optional(), enum: n(N()), default: N().optional() });
var iD = b({ type: _("string"), title: N().optional(), description: N().optional(), oneOf: n(b({ const: N(), title: N() })), default: N().optional() });
var nD = b({ type: _("string"), title: N().optional(), description: N().optional(), enum: n(N()), enumNames: n(N()).optional(), default: N().optional() });
var oD = G0([dD, iD]);
var rD = b({ type: _("array"), title: N().optional(), description: N().optional(), minItems: s().optional(), maxItems: s().optional(), items: b({ type: _("string"), enum: n(N()) }), default: n(N()).optional() });
var tD = b({ type: _("array"), title: N().optional(), description: N().optional(), minItems: s().optional(), maxItems: s().optional(), items: b({ anyOf: n(b({ const: N(), title: N() })) }), default: n(N()).optional() });
var aD = G0([rD, tD]);
var sD = G0([nD, oD, aD]);
var eD = G0([sD, lD, cD, pD]);
var Qw = B9.extend({ mode: _("form").optional(), message: N(), requestedSchema: b({ type: _("object"), properties: K0(N(), eD), required: n(N()).optional() }) });
var $w = B9.extend({ mode: _("url"), message: N(), elicitationId: N(), url: N().url() });
var Xw = G0([Qw, $w]);
var Yw = I0.extend({ method: _("elicitation/create"), params: Xw });
var Jw = a0.extend({ elicitationId: N() });
var Gw = s0.extend({ method: _("notifications/elicitation/complete"), params: Jw });
var D4 = E0.extend({ action: y0(["accept", "decline", "cancel"]), content: OY((Q) => Q === null ? void 0 : Q, K0(N(), G0([N(), s(), M0(), n(N())])).optional()) });
var Ww = b({ type: _("ref/resource"), uri: N() });
var Hw = b({ type: _("ref/prompt"), name: N() });
var Bw = p0.extend({ ref: G0([Hw, Ww]), argument: b({ name: N(), value: N() }), context: b({ arguments: K0(N(), N()).optional() }).optional() });
var r8 = I0.extend({ method: _("completion/complete"), params: Bw });
var zw = E0.extend({ completion: S0({ values: n(N()).max(100), total: L0(s().int()), hasMore: L0(M0()) }) });
var Kw = b({ uri: N().startsWith("file://"), name: N().optional(), _meta: K0(N(), z0()).optional() });
var Vw = I0.extend({ method: _("roots/list"), params: p0.optional() });
var vY = E0.extend({ roots: n(Kw) });
var qw = s0.extend({ method: _("notifications/roots/list_changed"), params: a0.optional() });
var xk = G0([T8, RY, r8, _Y, i8, d8, l8, c8, p8, VD, UD, O4, n8, y8, h8, f8, m8]);
var yk = G0([v8, x8, IY, qw, L9]);
var gk = G0([k8, O9, kY, D4, vY, g8, u8, F4]);
var hk = G0([T8, mD, Yw, Vw, y8, h8, f8, m8]);
var fk = G0([v8, x8, vD, FD, zD, SD, ED, L9, Gw]);
var uk = G0([k8, sO, zw, ID, DD, GD, WD, BD, o8, ZD, g8, u8, F4]);
var Fw = new Set("ABCDEFGHIJKLMNOPQRSTUVXYZabcdefghijklmnopqrstuvxyz0123456789");
var nK = G5(x7(), 1);
var oK = G5(iK(), 1);
var aK;
(function(Q) {
  Q.Completable = "McpCompletable";
})(aK || (aK = {}));
function YV(Q) {
  let $;
  return () => $ ??= Q();
}
var fP = YV(() => i1.object({ session_id: i1.string(), ws_url: i1.string(), work_dir: i1.string().optional(), session_key: i1.string().optional() }));
function Ph({ prompt: Q, options: $ }) {
  let { systemPrompt: X, settings: Y, settingSources: J, sandbox: G, ...W } = $ ?? {}, H, B;
  if (X === void 0) H = "";
  else if (typeof X === "string") H = X;
  else if (X.type === "preset") B = X.append;
  let z = W.pathToClaudeCodeExecutable;
  if (!z) {
    let x6 = lP(import.meta.url), y6 = HV(x6, "..");
    z = HV(y6, "cli.js");
  }
  process.env.CLAUDE_AGENT_SDK_VERSION = "0.2.81";
  let { abortController: K = g6(), additionalDirectories: U = [], agent: q, agents: V, allowedTools: L = [], betas: F, canUseTool: w, continue: D, cwd: M, debug: R, debugFile: Z, disallowedTools: v = [], tools: O0, env: D0, executable: d0 = h6() ? "bun" : "node", executableArgs: B6 = [], extraArgs: F1 = {}, fallbackModel: z6, enableFileCheckpointing: y1, toolConfig: K6, forkSession: h, hooks: S4, includePartialMessages: gQ, onElicitation: _4, persistSession: k4, thinking: V6, effort: c9, maxThinkingTokens: v6, maxTurns: C0, maxBudgetUsd: g1, mcpServers: T6, model: BV, outputFormat: e7, permissionMode: zV = "default", allowDangerouslySkipPermissions: KV = false, permissionPromptToolName: VV, plugins: qV, workload: Q5, resume: UV, resumeSessionAt: LV, sessionId: FV, stderr: NV, strictMcpConfig: OV } = W, $5 = e7?.type === "json_schema" ? e7.schema : void 0, q6 = D0;
  if (!q6) q6 = { ...process.env };
  if (!q6.CLAUDE_CODE_ENTRYPOINT) q6.CLAUDE_CODE_ENTRYPOINT = "sdk-ts";
  if (y1) q6.CLAUDE_CODE_ENABLE_SDK_FILE_CHECKPOINTING = "true";
  if (K6?.askUserQuestion?.previewFormat) q6.CLAUDE_CODE_QUESTION_PREVIEW_FORMAT = K6.askUserQuestion.previewFormat;
  if (!z) throw Error("pathToClaudeCodeExecutable is required");
  let hQ = {}, X5 = /* @__PURE__ */ new Map();
  if (T6) for (let [x6, y6] of Object.entries(T6)) if (y6.type === "sdk" && "instance" in y6) X5.set(x6, y6.instance), hQ[x6] = { type: "sdk", name: x6 };
  else hQ[x6] = y6;
  let DV = typeof Q === "string", v4;
  if (V6) switch (V6.type) {
    case "adaptive":
      v4 = { type: "adaptive" };
      break;
    case "enabled":
      v4 = { type: "enabled", budgetTokens: V6.budgetTokens };
      break;
    case "disabled":
      v4 = { type: "disabled" };
      break;
  }
  else if (v6 !== void 0) v4 = v6 === 0 ? { type: "disabled" } : { type: "enabled", budgetTokens: v6 };
  let Y5 = new y4({ abortController: K, additionalDirectories: U, agent: q, betas: F, cwd: M, debug: R, debugFile: Z, executable: d0, executableArgs: B6, extraArgs: Q5 ? { ...F1, workload: Q5 } : F1, pathToClaudeCodeExecutable: z, env: q6, forkSession: h, stderr: NV, thinkingConfig: v4, effort: c9, maxTurns: C0, maxBudgetUsd: g1, model: BV, fallbackModel: z6, jsonSchema: $5, permissionMode: zV, allowDangerouslySkipPermissions: KV, permissionPromptToolName: VV, continueConversation: D, resume: UV, resumeSessionAt: LV, sessionId: FV, settings: typeof Y === "object" ? W0(Y) : Y, settingSources: J ?? [], allowedTools: L, disallowedTools: v, tools: O0, mcpServers: hQ, strictMcpConfig: OV, canUseTool: !!w, hooks: !!S4, includePartialMessages: gQ, persistSession: k4, plugins: qV, sandbox: G, spawnClaudeCodeProcess: W.spawnClaudeCodeProcess }), wV = { systemPrompt: H, appendSystemPrompt: B, agents: V, promptSuggestions: W.promptSuggestions, agentProgressSummaries: W.agentProgressSummaries }, J5 = new h4(Y5, DV, w, S4, K, X5, $5, wV, _4);
  if (typeof Q === "string") Y5.write(W0({ type: "user", session_id: "", message: { role: "user", content: [{ type: "text", text: Q }] }, parent_tool_use_id: null }) + `
`);
  else J5.streamInput(Q);
  return J5;
}

// index.ts
function emit(obj) {
  process.stdout.write(JSON.stringify(obj) + "\n");
}
function log(msg) {
  process.stderr.write(`[bridge] ${msg}
`);
}
function buildSDKOptions(opts) {
  if (!opts) return {};
  const sdkOpts = {};
  if (opts.model) sdkOpts.model = opts.model;
  if (opts.cwd) sdkOpts.cwd = opts.cwd;
  if (opts.system_prompt) sdkOpts.systemPrompt = opts.system_prompt;
  if (opts.resume) sdkOpts.resume = opts.resume;
  if (opts.max_turns) sdkOpts.maxTurns = opts.max_turns;
  if (opts.permission_mode) {
    sdkOpts.permissionMode = opts.permission_mode;
    if (opts.permission_mode === "bypassPermissions") {
      sdkOpts.allowDangerouslySkipPermissions = true;
    }
  }
  if (opts.mcp_servers) sdkOpts.mcpServers = opts.mcp_servers;
  if (opts.allowed_tools) sdkOpts.allowedTools = opts.allowed_tools;
  if (opts.no_user_settings) {
    sdkOpts.settingSources = [];
  } else {
    sdkOpts.settingSources = ["user", "project", "local"];
  }
  if (opts.disabled_tools && opts.disabled_tools.length > 0) {
    sdkOpts.disallowedTools = opts.disabled_tools;
  }
  return sdkOpts;
}
function extractText(content) {
  if (!Array.isArray(content)) return "";
  return content.filter(
    (block) => typeof block === "object" && block !== null && "type" in block && block.type === "text" && "text" in block
  ).map((block) => block.text).join("");
}
async function handleQuery(req) {
  const reqId = req.request_id || "";
  const emitReq = (obj) => emit({ ...obj, request_id: reqId });
  const sdkOptions = buildSDKOptions(req.options);
  log(`query start \u2014 rid=${reqId} model=${sdkOptions.model ?? "default"} prompt="${req.prompt.slice(0, 80)}..."`);
  const timeoutMs = 10 * 60 * 1e3;
  const timeout = setTimeout(() => {
    log(`query timeout \u2014 rid=${reqId} no result after 10 minutes`);
    emitReq({ event: "error", message: "query timeout: no result after 10 minutes" });
  }, timeoutMs);
  try {
    const stream = Ph({
      prompt: req.prompt,
      options: sdkOptions
    });
    for await (const message of stream) {
      const msg = message;
      const msgType = msg.type;
      switch (msgType) {
        // ── System init ──────────────────────────────────────────────
        case "system": {
          emitReq({
            event: "system",
            session_id: msg.session_id,
            tools: msg.tools,
            model: msg.model
          });
          break;
        }
        // ── Assistant text + tool_use blocks ────────────────────────
        case "assistant": {
          const inner = msg.message;
          if (inner?.content && Array.isArray(inner.content)) {
            const text = extractText(inner.content);
            if (text) {
              emitReq({ event: "assistant", text });
            }
            for (const block of inner.content) {
              if (block.type === "tool_use") {
                emitReq({
                  event: "tool_use",
                  id: block.id,
                  name: block.name,
                  input: block.input
                });
              }
            }
          }
          break;
        }
        // ── Tool use summary ─────────────────────────────────────────
        case "tool_use_summary": {
          emitReq({
            event: "tool_result",
            content: msg.summary
          });
          break;
        }
        // ── Result (success or error) ────────────────────────────────
        case "result": {
          const subtype = msg.subtype;
          if (subtype === "success") {
            emitReq({
              event: "result",
              content: msg.result,
              cost_usd: msg.total_cost_usd,
              session_id: msg.session_id,
              duration_ms: msg.duration_ms,
              num_turns: msg.num_turns
            });
          } else {
            const errors = msg.errors;
            emitReq({
              event: "error",
              message: errors?.join("; ") ?? `result error: ${subtype}`,
              subtype: subtype ?? "unknown"
            });
          }
          break;
        }
        // ── All other message types (status, hooks, etc.) ───────────
        default: {
          break;
        }
      }
    }
  } catch (err) {
    const errMsg = err instanceof Error ? err.message : String(err);
    log(`query error: rid=${reqId} ${errMsg}`);
    emitReq({ event: "error", message: errMsg });
  } finally {
    clearTimeout(timeout);
  }
}
async function handleRequest(line) {
  let req;
  try {
    req = JSON.parse(line);
  } catch {
    emit({ event: "error", message: `invalid JSON: ${line.slice(0, 200)}` });
    return;
  }
  if (!req.command) {
    emit({ event: "error", request_id: req.request_id || "", message: "missing 'command' field" });
    return;
  }
  const reqId = req.request_id || "";
  switch (req.command) {
    case "query": {
      if (!req.prompt) {
        emit({ event: "error", request_id: reqId, message: "missing 'prompt' field for query command" });
        return;
      }
      await handleQuery(req);
      break;
    }
    case "ping": {
      emit({ event: "pong", request_id: reqId });
      break;
    }
    default: {
      emit({ event: "error", request_id: reqId, message: `unknown command: ${req.command}` });
    }
  }
}
function main() {
  log("bridge started \u2014 waiting for commands on stdin");
  const rl = createInterface({
    input: process.stdin,
    terminal: false
  });
  rl.on("line", (line) => {
    const trimmed = line.trim();
    if (!trimmed) return;
    handleRequest(trimmed).catch((err) => {
      const errMsg = err instanceof Error ? err.message : String(err);
      log(`unhandled error in request processing: ${errMsg}`);
      emit({ event: "error", message: `internal bridge error: ${errMsg}` });
    });
  });
  rl.on("close", () => {
    log("stdin closed \u2014 shutting down");
    process.exit(0);
  });
  process.on("unhandledRejection", (reason) => {
    const msg = reason instanceof Error ? reason.message : String(reason);
    log(`unhandled rejection: ${msg}`);
    emit({ event: "error", message: `unhandled rejection: ${msg}` });
  });
  process.on("uncaughtException", (err) => {
    log(`uncaught exception: ${err.message}`);
    emit({ event: "error", message: `uncaught exception: ${err.message}` });
    process.exit(1);
  });
}
main();
