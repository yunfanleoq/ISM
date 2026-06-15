<template>
  <div >
    <DataGrid ref="dg" :data="loggerList"  :pagination="true"  :pageSize="pageSize" :rowCss="getRowCss" style="">
      <GridColumn field="id" width="5%" :title="$t('displayConfig.Logger.ID')"></GridColumn>
      <GridColumn field="level" width="10%" :title="$t('displayConfig.Logger.Level')"></GridColumn>
      <GridColumn field="content" width="50%" :title="$t('displayConfig.Logger.Content')">
      </GridColumn>
      <GridColumn field="time" :sortable="true" width="10%" :title="$t('displayConfig.Logger.Time')">
      </GridColumn>
    </DataGrid>
  </div>
</template>

<script>
import { mapActions, mapGetters, mapState, mapMutations } from 'vuex'
import store from "../../store";

export default {
  name: "ISMLogger",
  i18n: require('../../i18n/language'),
  components: {

  },
  data() {
    return {
      pageSize: 10,
      operators: ["nofilter", "equal", "notequal", "less", "greater"],
      status: [
        { value: null, text: "All" },
        { value: "P", text: "P" },
        { value: "N", text: "N" }
      ]
    }
  },
  props: [],
  computed: {
    ...mapState({
      loggerList: state => store.state.ISMDisPlayEditorTool.loggerList,
    }),
  },
  methods: {
    getRowCss(row) {
      if (row.level == "debug") {
        return { background: "#ffff"};
      }
      else if (row.level == "info") {
        return { background: "#b8eecf" };
      }
      else if (row.level == "warn") {
        return { background: "#faad14"};
      }
      else if (row.level == "error") {
        return { background: "#f5222d" };
      }
      return null;
    },
  },
  created(){

  },
  mounted() {

  }
}
</script>

<style scoped>

</style>
