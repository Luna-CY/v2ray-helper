<template>
  <el-dialog :model-value="show" destroy-on-close width="80%" @close="close">
    <el-table :data="data" border style="width: 100%">
      <el-table-column label="系统类型" prop="os"></el-table-column>
      <el-table-column label="客户端" prop="client"></el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-link :href="scope.row.git" v-if="'#' !== scope.row.git" :underline="false" class="margin-right"
                   target="_blank" type="primary">GitHub
          </el-link>
          <el-link :href="scope.row.address" :underline="false" class="margin-right" target="_blank" type="primary">下载
          </el-link>
          <el-link @click="$router.push({name: scope.row.helper})" v-if="'#' !== scope.row.helper" :underline="false"
                   type="primary">配置帮助
          </el-link>
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"

export default defineComponent({
  name: "Download",

  props: {
    show: Boolean
  },

  emits: ['update:show'],

  data() {
    return {
      data: [
        {
          os: "Windows",
          client: "V2rayN",
          git: "https://github.com/2dust/v2rayN/releases",
          address: "https://github.com/2dust/v2rayN/releases/download/5.23/v2rayN-Core.zip",
          helper: "HelperV2rayN",
        },
        {
          os: "MacOSX",
          client: "V2rayX",
          git: "https://github.com/Cenmrev/V2RayX/releases",
          address: "https://github.com/Cenmrev/V2RayX/releases/download/v1.5.1/V2RayX.app.zip",
          helper: "HelperV2rayX",
        },
        {
          os: "Android",
          client: "SagerNet",
          git: "https://github.com/SagerNet/SagerNet/releases",
          address: "https://github.com/SagerNet/SagerNet/releases/download/0.8-beta06/SN-0.8-beta06-arm64-v8a.apk",
          helper: "HelperSagerNet",
        },
        {
          os: "IOS",
          client: "i2Ray",
          git: "#",
          address: 'https://apps.apple.com/us/app/i2ray/id1445270056',
          helper: "#",
        },
      ],
    }
  },

  methods: {
    close() {
      this.$emit('update:show', false)
    },
  },
})
</script>

<style lang="scss" scoped>
.download-client-box-md-and-up {
  height: 300px;
  width: 100%;
  display: flex;
  justify-content: space-around;
  align-items: center;

  .download-client-item {
    width: 100%;
    display: flex;
    flex-flow: column;
    align-items: center;

    .client-image {
      width: 128px;
      height: 128px;
    }
  }
}

.download-client-box-sm-and-down {
  display: flex;
  flex-flow: column;
  align-items: center;

  .download-client-item {
    width: 100%;
    display: flex;
    flex-flow: column;
    align-items: center;

    .client-image {
      width: 64px;
      height: 64px;
    }
  }
}
</style>
