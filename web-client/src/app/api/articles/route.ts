import { Pool } from "pg"
import { NextResponse } from "next/server"

// データベース接続プールを作成
const pool = new Pool({
  connectionString: process.env.DB_URL,
})

export async function GET() {
  try {
    // 記事を取得するクエリ（RSSフィードの情報も一緒に取得）
    const query = `
      SELECT 
        a.id, 
        a.rss_id, 
        a.title, 
        a.link, 
        a.description, 
        a.published_at,
        a.summary
      FROM 
        articles a
      ORDER BY 
        a.published_at DESC
      LIMIT 50
    `

    const client = await pool.connect()
    const result = await client.query(query)
    client.release()

    return NextResponse.json({
      success: true,
      articles: result.rows,
    })
  } catch (error) {
    console.error("記事の取得中にエラーが発生しました:", error)
    return NextResponse.json(
      {
        success: false,
        error: "記事の取得に失敗しました",
      },
      { status: 500 }
    )
  }
}
