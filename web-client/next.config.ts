/** @type {import('next').NextConfig} */
import * as dotenv from "dotenv"
import path from "path"

// ルートディレクトリの.envファイルを読み込む
dotenv.config({ path: path.resolve(__dirname, "../.env") })

const nextConfig = {
  // 既存の設定を保持
}

export default nextConfig
