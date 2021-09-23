<template>
  <el-dialog :model-value="show" destroy-on-close width="80%" @close="close">
    <el-form :model="form" label-width="120px">
      <div class="inline-form-item-2">
        <el-form-item class="form-item" label="域名/IP" prop="host">
          <el-input v-model="form.host" readonly></el-input>
        </el-form-item>
        <el-form-item class="form-item" label="端口" label-width="60px" prop="port">
          <el-input v-model="form.port" readonly></el-input>
        </el-form-item>
      </div>
      <div class="inline-form-item-2">
        <el-form-item class="form-item" label="用户ID" prop="user_id">
          <el-input v-model="form.user_id" readonly></el-input>
        </el-form-item>
        <el-form-item class="form-item" label="额外ID" label-width="60px" prop="alter_id">
          <el-input v-model="form.alter_id" readonly></el-input>
        </el-form-item>
      </div>
      <div class="inline-form-item-2">
        <el-form-item class="form-item-0" label="使用TLS">
          <el-switch v-model="form.use_tls" disabled></el-switch>
        </el-form-item>
        <el-form-item class="form-item-1" label="传输方式" label-width="80px" prop="transport_type">
          <el-select v-model="form.transport_type" class="w-100" disabled>
            <el-option :value="1" label="TCP"></el-option>
            <el-option :value="2" label="WebSocket"></el-option>
            <el-option :value="3" label="KCP"></el-option>
            <el-option :value="4" label="HTTP2"></el-option>
          </el-select>
        </el-form-item>
      </div>
      <template v-if="1 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">TCP传输配置</el-divider>
        <el-form-item label="伪装类型">
          <el-select v-model="form.tcp.type" disabled>
            <el-option label="NONE" value="none"></el-option>
            <el-option label="HTTP" value="http"></el-option>
          </el-select>
        </el-form-item>
        <template v-if="'http' === form.tcp.type">
          <div class="inline-form-item-2">
            <el-form-item class="form-item" label="请求版本">
              <el-select v-model="form.tcp.request.version" disabled>
                <el-option label="HTTP/1.1" value="1.1"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="form-item" label="请求方式">
              <el-select v-model="form.tcp.request.method" disabled>
                <el-option label="GET" value="GET"></el-option>
                <el-option label="POST" value="POST"></el-option>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item label="请求路径">
            <el-input v-model="form.tcp.request.path" readonly></el-input>
          </el-form-item>
          <template v-for="header in form.tcp.request.headers" v-bind:key="header">
            <div class="inline-form-item-2">
              <el-form-item class="form-item" label="字段名">
                <el-input v-model="header.key" readonly></el-input>
              </el-form-item>
              <el-form-item class="form-item" label="字段值">
                <el-input v-model="header.value" readonly></el-input>
              </el-form-item>
            </div>
          </template>
          <div class="inline-form-item-3">
            <el-form-item class="form-item" label="响应版本">
              <el-select v-model="form.tcp.response.version" disabled>
                <el-option label="HTTP/1.1" value="1.1"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="form-item" label="状态码">
              <el-input v-model="form.tcp.response.status" readonly></el-input>
            </el-form-item>
            <el-form-item class="form-item" label="状态描述">
              <el-input v-model="form.tcp.response.reason" readonly></el-input>
            </el-form-item>
          </div>
          <template v-for="header in form.tcp.response.headers" v-bind:key="header">
            <div class="inline-form-item-2">
              <el-form-item class="form-item" label="字段名">
                <el-input v-model="header.key" readonly></el-input>
              </el-form-item>
              <el-form-item class="form-item" label="字段值">
                <el-input v-model="header.value" readonly></el-input>
              </el-form-item>
            </div>
          </template>
        </template>
      </template>
      <template v-if="2 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">WebSocket传输配置</el-divider>
        <el-form-item label="路径">
          <el-input v-model="form.web_socket.path" readonly></el-input>
        </el-form-item>
        <template v-for="header in form.web_socket.headers" v-bind:key="header">
          <div class="inline-form-item-2">
            <el-form-item class="form-item" label="字段名">
              <el-input v-model="header.key" readonly></el-input>
            </el-form-item>
            <el-form-item class="form-item" label="字段值">
              <el-input v-model="header.value" readonly></el-input>
            </el-form-item>
          </div>
        </template>
      </template>
      <template v-if="3 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">KCP传输配置</el-divider>
        <div class="inline-form-item-2">
          <el-form-item class="form-item" label="伪装类型">
            <el-select v-model="form.kcp.type" class="w-100" disabled>
              <el-option label="none" value="none"></el-option>
              <el-option label="srtp" value="srtp"></el-option>
              <el-option label="utp" value="utp"></el-option>
              <el-option label="wechat-video" value="wechat-video"></el-option>
              <el-option label="dtls" value="dtls"></el-option>
              <el-option label="wireguard" value="wireguard"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item class="form-item" label="开启拥塞控制">
            <el-switch v-model="form.kcp.congestion" disabled></el-switch>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item class="form-item" label="MTU大小">
            <el-input v-model="form.kcp.mtu" readonly></el-input>
          </el-form-item>
          <el-form-item class="form-item" label="TTI间隔时间">
            <el-input v-model="form.kcp.tti" readonly></el-input>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item class="form-item" label="上行带宽">
            <el-input v-model="form.kcp.uplink_capacity" readonly>
              <template #append>MB/S</template>
            </el-input>
          </el-form-item>
          <el-form-item class="form-item" label="读取缓冲区大小">
            <el-input v-model="form.kcp.read_buffer_size" readonly>
              <template #append>MB</template>
            </el-input>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item class="form-item" label="下行带宽">
            <el-input v-model="form.kcp.downlink_capacity" readonly>
              <template #append>MB/S</template>
            </el-input>
          </el-form-item>
          <el-form-item class="form-item" label="写入缓冲区大小">
            <el-input v-model="form.kcp.write_buffer_size" readonly>
              <template #append>MB</template>
            </el-input>
          </el-form-item>
        </div>
      </template>
      <template v-if="4 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">HTTP2传输配置</el-divider>
        <el-form-item label="域名">
          <el-input v-model="form.http2.host" readonly></el-input>
        </el-form-item>
        <el-form-item label="路径">
          <el-input v-model="form.http2.path" readonly></el-input>
        </el-form-item>
      </template>
    </el-form>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import {V2rayEndpointDetailData} from "@/api/v2ray_endpoint_detail"

export default defineComponent({
  name: "EndpointDetail",

  props: {
    show: Boolean,
    data: Object,
  },

  emits: ['update:show'],

  watch: {
    show: function () {
      if (this.show) {
        this.form = this.data as V2rayEndpointDetailData
      }
    }
  },

  data() {
    return {
      form: new V2rayEndpointDetailData(),
    }
  },

  methods: {
    close() {
      this.$emit('update:show', false)
      this.form = new V2rayEndpointDetailData()
    },
  },
})
</script>

<style lang="scss" scoped>
.inline-form-item-2 {
  display: flex;

  .form-item {
    width: 50%;
  }

  .form-item-0 {
    flex-grow: 0;
  }

  .form-item-1 {
    flex-grow: 1;
  }
}
</style>
