<template>
  <div class="content-center el-container">
    <div class="el-main el-main-md-and-up hidden-sm-and-down">
      <div class="margin-bottom-2x">
        <div class="home-header">
          <div class="left-notice">
            <el-alert type="error" effect="dark" :closable="false"
                      v-if="$store.getters.local.is_default_key || $store.getters.local.is_default_key">
              <template #title>
                <span v-if="$store.getters.local.is_default_key">当前应用使用的口令为默认口令，请及时修改以免造成损失</span>
                <span v-if="$store.getters.local.is_default_key && $store.getters.local.is_default_remove_key">; </span>
                <span v-if="$store.getters.local.is_default_remove_key">当前应用使用的配置删除口令为默认口令，请及时修改以免造成损失</span>
              </template>
            </el-alert>
          </div>
          <div class="right-buttons margin-left">
            <i class="el-icon-bell" @click="$message.info('暂未实现此功能')"></i>
            <i class="el-icon-setting margin-left" @click="$message.info('暂未实现此功能')"></i>
          </div>
        </div>
        <el-divider></el-divider>
      </div>
      <template v-for="item in data" v-bind:key="item.id">
        <div class="endpoint-box el-row margin-bottom-2x" v-on:mouseout="item.show_delete_button = false"
             v-on:mouseover="item.show_delete_button = true">
          <div class="el-col-4">{{ item.host }}</div>
          <div class="el-col-2">{{ item.port }}</div>
          <div class="el-col-4">{{ item.user_id }}</div>
          <div class="el-col-2">{{ item.alter_id }}</div>
          <div class="el-col-2">{{ getTransportType(item.transport_type) }}</div>
          <div class="el-col-xl-6 el-col-lg-4 el-col-md-4">{{ item.remark ? item.remark : '-' }}</div>
          <div class="el-col-xl-4 el-col-lg-6 el-col-md-6">
            <el-button :loading="item.downloading" type="primary" @click="download(item)">生成VMess链接</el-button>
            <el-button :loading="item.loading" class="margin-left" type="primary" @click="showDetail(item)">显示完整配置
            </el-button>
            <el-button v-show="item.show_delete_button" circle class="delete-button" icon="el-icon-delete"
                       type="danger" @click="removeItem = item; showRemoveModal = true"></el-button>
          </div>
        </div>
      </template>

      <div v-if="0 === data.length" class="endpoint-box el-row margin-bottom-2x">
        <div class="el-col-24">暂时没有可用的配置列表</div>
      </div>

      <div class="el-row">
        <div class="el-col-24">
          <el-button size="medium" type="success" @click="showNewModal = true">添加配置</el-button>
          <el-button size="medium" type="primary" @click="showDevelopV2rayModal = true">部署服务器</el-button>
          <el-button size="medium" type="primary" @click="showDownloadModal = true">下载客户端</el-button>
        </div>
      </div>
    </div>

    <div class="el-main el-main-sm-and-down hidden-md-and-up">
      <div class="endpoint-box el-row margin-bottom">
        <div class="el-col-24">暂未适配手机端</div>
      </div>
    </div>
  </div>

  <Remove v-model:show="showRemoveModal" v-bind:loading="removing" v-on:confirm="remove"></Remove>
  <NewV2rayEndpoint v-model:show="showNewModal" v-on:success="load"></NewV2rayEndpoint>
  <Download v-model:show="showDownloadModal"></Download>
  <QRCode v-model:show="showQRCodeModal" v-bind:content="v2rayNgVMessString"></QRCode>
  <DeployV2rayServer v-model:show="showDevelopV2rayModal" v-on:success="load"></DeployV2rayServer>
  <EndpointDetail v-model:show="showDetailModal" v-bind:data="endpointDetail"></EndpointDetail>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import axios, {AxiosResponse} from "axios"
import {API_V2RAY_ENDPOINT_LIST, V2rayEndpointListItem, V2rayEndpointListResponse} from "@/api/v2ray_endpoint_list"
import {
  API_V2RAY_ENDPOINT_DOWNLOAD,
  V2rayEndpointDownloadForm,
  V2rayEndpointDownloadResponse
} from "@/api/v2ray_endpoint_download"
import {API_V2RAY_ENDPOINT_REMOVE, V2rayEndpointRemoveForm} from "@/api/v2ray_endpoint_remove"
import {BaseResponse} from "@/api/base"
import Remove from "@/components/Remove.vue"
import NewV2rayEndpoint from "@/components/NewV2rayEndpoint.vue"
import Download from "@/components/Download.vue"
import QRCode from "@/components/QRCode.vue"
import DeployV2rayServer from "@/components/DeployV2rayServer.vue"
import {
  API_V2RAY_ENDPOINT_DETAIL,
  V2rayEndpointDetailData,
  V2rayEndpointDetailParams,
  V2rayEndpointDetailResponse
} from "@/api/v2ray_endpoint_detail"
import EndpointDetail from "@/components/EndpointDetail.vue";
import {API_META_INFO, MetaInfoResponse} from "@/api/meta_info";
import {StoryStateLocal} from "@/store";

const md5 = require('md5')

export default defineComponent({
  name: 'Home',

  components: {EndpointDetail, DeployV2rayServer, QRCode, Download, NewV2rayEndpoint, Remove},

  data() {
    return {
      loading: true,
      data: new Array<V2rayEndpointListItem>(),
      showRemoveModal: false,
      removeItem: new V2rayEndpointListItem(),
      removing: false,
      showNewModal: false,
      showDownloadModal: false,
      showQRCodeModal: false,
      showDevelopV2rayModal: false,
      v2rayNgVMessString: "",
      endpointDetail: new V2rayEndpointDetailData(),
      showDetailModal: false,
    }
  },

  mounted() {
    axios.get(API_META_INFO).then((response: AxiosResponse<MetaInfoResponse>) => {
      if (0 != response.data.code) {
        this.$message.error(response.data.message)

        return
      }

      let state = new StoryStateLocal()
      state.is_default_key = response.data.data.is_default_key
      state.is_default_remove_key = response.data.data.is_default_remove_key

      this.$store.commit('local', state)
    })

    this.load()
  },

  methods: {
    load() {
      this.loading = true

      axios.get(API_V2RAY_ENDPOINT_LIST).then((response: AxiosResponse<V2rayEndpointListResponse>) => {
        this.loading = false

        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        this.data = response.data.data.data
      })
    },
    getTransportType(transportType: number) {
      let map = {1: "TCP", 2: "WebSocket", 3: "KCP", 4: "HTTP2"} as any

      return map[transportType]
    },
    download(item: V2rayEndpointListItem) {
      item.downloading = true

      let form = new V2rayEndpointDownloadForm()
      form.id = item.id

      axios.post(API_V2RAY_ENDPOINT_DOWNLOAD, form).then((response: AxiosResponse<V2rayEndpointDownloadResponse>) => {
        item.downloading = false

        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        this.v2rayNgVMessString = response.data.data.content
        this.showQRCodeModal = true
      })
    },
    remove(password: string) {
      const body = new V2rayEndpointRemoveForm()
      body.id = this.removeItem.id
      body.password = md5(password)

      this.removing = true
      axios.post(API_V2RAY_ENDPOINT_REMOVE, body).then((response: AxiosResponse<BaseResponse>) => {
        this.removing = false

        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        this.showRemoveModal = false
        this.$message.success("删除成功")

        this.load()
      })
    },

    showDetail(item: V2rayEndpointListItem) {
      item.loading = true

      let params = new V2rayEndpointDetailParams()
      params.id = item.id

      axios.get(API_V2RAY_ENDPOINT_DETAIL, {params}).then((response: AxiosResponse<V2rayEndpointDetailResponse>) => {
        item.loading = false

        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        this.endpointDetail = response.data.data
        this.showDetailModal = true
      })
    },
  },
});
</script>

<style lang="scss">
body {
  background-color: #2d3a4b;
}
</style>

<style lang="scss">
.el-row {
  align-items: center;
}

.el-main.el-main-md-and-up {
  --el-main-padding: 30px;

  .endpoint-box {
    background-color: #fff;
    padding: 25px;

    .delete-button {
      position: absolute;
      top: -15px;
      right: -15px;
    }
  }
}

.el-main.el-main-sm-and-down {
  --el-main-padding: 15px;

  .endpoint-box {
    background-color: #fff;
    padding: 15px;
  }
}

.home-header {
  color: #FFFFFF;
  font-size: 120%;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .left-notice {
    flex-grow: 1;
  }

  i {
    cursor: pointer;
  }
}
</style>
