<template>
  <div>
    <div class="sp-container">
      <sp-sign-in />
      <sp-bank-balances />
      <sp-token-send />
      <!-- this line is used by starport scaffolding # 4 -->
		<sp-type-form path="jack139.contract.contract" type="contract" :fields="[ ['creator', 1, 'string'] , ['contractNo', 2, 'string'] , ['partyA', 3, 'string'] , ['partyB', 4, 'string'] , ['action', 5, 'string'] , ['data', 6, 'string'] ]" />
    </div>
  </div>
</template>

<script>
import * as sp from "@tendermint/vue";

export default {
  components: { ...sp },
};
</script>
