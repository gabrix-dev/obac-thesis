import { reactive } from "vue";

export const store = reactive({
  rules: [],
  addDefaultRule(ruleName) {
    const firstConstraint = {
      key: "EventID",
      value: ["Qatar2022"],
      operator: "Equals",
    };
    const rule = {
      name: ruleName,
      action: "Permit",
      ownershipConstraints: [
        {
          smartContract: {
            address: "0xfd6b508D7A6253809D796101b2B6ae9164ab5959",
            blockchain: "ethereum",
          },
          metadataConstraints: [firstConstraint],
        },
      ],
    };
    this.rules.push(rule);
  },
  updateRule(index, rule) {
    this.rules[index] = rule;
  },
  getRules() {
    return this.rules;
  },
  deleteRule(index) {
    this.rules.splice(index, 1);
  },
  setRules(backendRules) {
    this.rules = backendRules;
  },
  scopes: [],
  setScopes(scopes) {
    this.scopes = scopes;
  },
  getScopes() {
    return this.scopes;
  },
  getScope(index) {
    return this.scopes[index];
  },
  getName(index) {
    return this.scopes[index].name;
  },
  setClaimType(value, claimIndex, scopeIndex) {
    this.scopes[scopeIndex].claims[claimIndex].type = value;
  },
  setClaimValue(value, claimIndex, scopeIndex) {
    this.scopes[scopeIndex].claims[claimIndex].value = value;
  },
  setClaimOidc(value, claimIndex, scopeIndex) {
    this.scopes[scopeIndex].claims[claimIndex].oidc = value;
  },
  addScope(name) {
    const scope = {
      name: name,
      claims: [
        {
          type: "",
          value: "",
          oidc: "",
        },
      ],
    };
    this.scopes.push(scope);
  },
  addDefaultClaim(scopeIndex) {
    const claim = {
      type: "",
      value: "",
      oidc: "",
    };
    this.scopes[scopeIndex].claims.push(claim);
  },
  setScopeName(scopeIndex, name) {
    this.scopes[scopeIndex].name = name;
  },
  deleteClaim(scopeIndex, claimIndex) {
    if (this.scopes[scopeIndex].claims.length != 1) {
      this.scopes[scopeIndex].claims.splice(claimIndex, 1);
    }
  },
  deleteScope(scopeIndex) {
    this.scopes.splice(scopeIndex, 1);
  },
});
