<template>
  <el-dialog :close-on-click-modal="false" :close-on-press-escape="false" :model-value="show" destroy-on-close
             width="80%" @close="close" :fullscreen="true">
    <el-form ref="NewV2rayEndpointForm" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="备注信息">
        <el-input v-model="form.remark" placeholder="备注信息"></el-input>
      </el-form-item>
      <div class="inline-form-item-2">
        <el-form-item class="form-item" label="域名/IP" prop="host">
          <el-input v-model="form.host" placeholder="主机地址: IP或域名"></el-input>
        </el-form-item>
        <el-form-item class="form-item" label="端口" label-width="60px" prop="port">
          <el-input v-model="form.port" max="65535" min="1" placeholder="端口" type="number"></el-input>
        </el-form-item>
      </div>
      <div class="inline-form-item-2">
        <el-form-item class="form-item" label="用户ID" prop="user_id">
          <el-input v-model="form.user_id" placeholder="用户身份ID"></el-input>
        </el-form-item>
        <el-form-item class="form-item" label="额外ID" label-width="60px" prop="alter_id">
          <el-input v-model="form.alter_id" max="65535" min="0" placeholder="额外ID的数量，需要与服务器配置一致"
                    type="number"></el-input>
        </el-form-item>
      </div>
      <div class="inline-form-item-2">
        <el-form-item class="form-item-0" label="使用TLS">
          <el-switch v-model="form.use_tls" :disabled="3 === parseInt(form.transport_type.toString())"></el-switch>
        </el-form-item>
        <el-form-item class="form-item-1" label="传输方式" label-width="80px" prop="transport_type">
          <el-select v-model="form.transport_type" class="w-100">
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
          <el-select v-model="form.tcp.type">
            <el-option label="NONE" value="none"></el-option>
            <el-option label="HTTP" value="http"></el-option>
          </el-select>
        </el-form-item>
        <template v-if="'http' === form.tcp.type">
          <div class="inline-form-item-2">
            <el-form-item class="form-item" label="请求版本">
              <el-select v-model="form.tcp.request.version">
                <el-option label="HTTP/1.1" value="1.1"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="form-item" label="请求方式">
              <el-select v-model="form.tcp.request.method">
                <el-option label="GET" value="GET"></el-option>
                <el-option label="POST" value="POST"></el-option>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item label="请求路径">
            <el-input v-model="form.tcp.request.path" placeholder="请求路径，多个路径用英文,分隔，每次请求会随机选择一个，默认为/"></el-input>
          </el-form-item>
          <template v-for="header in form.tcp.request.headers" v-bind:key="header">
            <div class="inline-form-item-2">
              <el-form-item class="form-item" label="字段名">
                <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
              </el-form-item>
              <el-form-item class="form-item" label="字段值">
                <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
              </el-form-item>
            </div>
          </template>
          <el-form-item class="content-center" label-width="0">
            <el-button type="primary" @click="addTcpRequestHeader">添加自定义请求头字段</el-button>
          </el-form-item>
          <div class="inline-form-item-3">
            <el-form-item class="form-item" label="响应版本">
              <el-select v-model="form.tcp.response.version">
                <el-option label="HTTP/1.1" value="1.1"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item class="form-item" label="状态码">
              <el-input v-model="form.tcp.response.status" placeholder="HTTP响应的状态码，默认为200"></el-input>
            </el-form-item>
            <el-form-item class="form-item" label="状态描述">
              <el-input v-model="form.tcp.response.reason" placeholder="HTTP响应的状态描述，默认为OK"></el-input>
            </el-form-item>
          </div>
          <template v-for="header in form.tcp.response.headers" v-bind:key="header">
            <div class="inline-form-item-2">
              <el-form-item class="form-item" label="字段名">
                <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
              </el-form-item>
              <el-form-item class="form-item" label="字段值">
                <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
              </el-form-item>
            </div>
          </template>
          <el-form-item class="content-center" label-width="0">
            <el-button type="primary" @click="addTcpResponseHeader">添加自定义响应头字段</el-button>
          </el-form-item>
        </template>
      </template>
      <template v-if="2 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">WebSocket传输配置</el-divider>
        <el-form-item label="路径">
          <el-input v-model="form.web_socket.path" placeholder="URI路径"></el-input>
        </el-form-item>
        <template v-for="header in form.web_socket.headers" v-bind:key="header">
          <div class="inline-form-item-2">
            <el-form-item class="form-item" label="字段名">
              <el-input v-model="header.key" placeholder="自定义头的字段名称"></el-input>
            </el-form-item>
            <el-form-item class="form-item" label="字段值">
              <el-input v-model="header.value" placeholder="自定义头的字段值"></el-input>
            </el-form-item>
          </div>
        </template>
        <el-form-item class="content-center" label-width="0">
          <el-button type="primary" @click="addWebSocketHeader">添加自定义头字段</el-button>
        </el-form-item>
      </template>
      <template v-if="3 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">KCP传输配置</el-divider>
        <div class="inline-form-item-2">
          <el-form-item class="form-item" label="伪装类型">
            <el-select v-model="form.kcp.type" class="w-100">
              <el-option label="none" value="none"></el-option>
              <el-option label="srtp" value="srtp"></el-option>
              <el-option label="utp" value="utp"></el-option>
              <el-option label="wechat-video" value="wechat-video"></el-option>
              <el-option label="dtls" value="dtls"></el-option>
              <el-option label="wireguard" value="wireguard"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item class="form-item" label="开启拥塞控制">
            <el-switch v-model="form.kcp.congestion"></el-switch>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item class="form-item" label="MTU大小">
            <el-input v-model="form.kcp.mtu" max="1460" min="576" placeholder="传输单元大小，576-1460之间的整数，默认为1350"
                      type="number"></el-input>
          </el-form-item>
          <el-form-item class="form-item" label="TTI间隔时间">
            <el-input v-model="form.kcp.tti" max="100" min="10" placeholder="传输间隔时间，10-100之间的整数，默认为50"
                      type="number"></el-input>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item class="form-item" label="上行带宽">
            <el-input v-model="form.kcp.uplink_capacity" min="0" placeholder="上行带宽大小，默认为5，单位MB/s" type="number">
              <template #append>MB/S</template>
            </el-input>
          </el-form-item>
          <el-form-item class="form-item" label="读取缓冲区大小">
            <el-input v-model="form.kcp.read_buffer_size" min="1" placeholder="读取缓冲区大小，默认为2，单位MB" type="number">
              <template #append>MB</template>
            </el-input>
          </el-form-item>
        </div>
        <div class="inline-form-item-2">
          <el-form-item class="form-item" label="下行带宽">
            <el-input v-model="form.kcp.downlink_capacity" min="0" placeholder="下行带宽大小，默认为20，单位MB/s"
                      type="number">
              <template #append>MB/S</template>
            </el-input>
          </el-form-item>
          <el-form-item class="form-item" label="写入缓冲区大小">
            <el-input v-model="form.kcp.write_buffer_size" min="1" placeholder="写入缓冲区大小，默认为2，单位MB" type="number">
              <template #append>MB</template>
            </el-input>
          </el-form-item>
        </div>
      </template>
      <template v-if="4 === parseInt(form.transport_type.toString())">
        <el-divider content-position="left">HTTP2传输配置</el-divider>
        <el-form-item label="域名">
          <el-input v-model="form.http2.host" placeholder="HTTP2的域名，多个使用英文,分隔"></el-input>
        </el-form-item>
        <el-form-item label="路径">
          <el-input v-model="form.http2.path" placeholder="URI路径"></el-input>
        </el-form-item>
      </template>
      <el-form-item class="content-center" label-width="0">
        <el-button type="danger" @click="close">取消</el-button>
        <el-button :loading="saving" type="primary" @click="save">保存</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import axios, {AxiosResponse} from "axios"
import {BaseResponse, Header} from "@/api/base"
import {API_V2RAY_ENDPOINT_NEW, V2rayEndpointNewForm} from "@/api/v2ray_endpoint_new"

export default defineComponent({
  name: "NewV2rayEndpoint",

  props: {
    show: Boolean,
  },

  emits: ['update:show', 'success'],

  watch: {
    'form.transport_type': function () {
      if (3 == this.form.transport_type) {
        this.form.use_tls = false
      }
    }
  },

  data() {
    return {
      saving: false,
      form: new V2rayEndpointNewForm(),
      rules: {
        host: [{required: true, message: '必须填写主机地址', trigger: 'blur'}],
        port: [{required: true, message: '必须填写端口号', trigger: 'blur'}],
        user_id: [{required: true, message: '必须填写用户ID', trigger: 'blur'}],
      },
    }
  },

  methods: {
    close() {
      this.form = new V2rayEndpointNewForm()

      this.$emit('update:show', false)
    },

    success() {
      this.close()

      this.$emit('success')
    },

    save() {
      let form = this.$refs.NewV2rayEndpointForm as any

      form.validate((isValid: boolean) => {
        if (!isValid) {
          return
        }

        this.form.port = parseInt(this.form.port.toString())
        this.form.alter_id = parseInt(this.form.alter_id.toString())

        this.saving = true
        axios.post(API_V2RAY_ENDPOINT_NEW, this.form).then((response: AxiosResponse<BaseResponse>) => {
          this.saving = false
          if (0 != response.data.code) {
            this.$message.error(response.data.message)

            return
          }

          this.$message.success('添加配置成功')
          this.success()
        })
      })
    },

    addWebSocketHeader() {
      this.form.web_socket.headers.push(new Header())
    },

    addTcpRequestHeader() {
      this.form.tcp.request.headers.push(new Header())
    },

    addTcpResponseHeader() {
      this.form.tcp.response.headers.push(new Header())
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
