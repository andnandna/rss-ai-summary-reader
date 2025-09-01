"use client";

import { useState, useEffect } from 'react';
import { FaRss, FaClock, faExternalLinkAlt } from 'react-icons/fa';

interface Article {
  id: number;
  rss_id: number;
  title: string;
  link: string;
  description: string;
  summary?: string;
  published_at: string;
}

export default function Home() {
  const [articles, setArticles] = useState<Article[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchArticles() {
      try {
        const response = await fetch('/api/articles');
        const data = await response.json();

        if (data.success) {
          setArticles(data.articles);
        } else {
          setError(data.error || '記事の取得に失敗しました');
        }
      } catch (err) {
        setError('ネットワークエラーが発生しました');
      } finally {
        setIsLoading(false);
      }
    }

    fetchArticles();
  }, []);

  if (isLoading) return (
    <div className="flex justify-center items-center min-h-screen bg-gray-100">
      <div className="text-center">
        <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-blue-500 mx-auto mb-4"></div>
        <p className="text-gray-600">読み込み中...</p>
      </div>
    </div>
  );

  if (error) return (
    <div className="flex justify-center items-center min-h-screen bg-red-50">
      <div className="bg-white p-8 rounded-lg shadow-md text-center">
        <h2 className="text-2xl text-red-600 mb-4">エラー</h2>
        <p className="text-gray-700">{error}</p>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        <header className="mb-10 text-center">
          <h1 className="text-4xl font-extrabold text-gray-900 mb-4">
            最新のRSSフィード
          </h1>
          <p className="text-xl text-gray-600">
            さまざまなソースから集められた最新の記事
          </p>
        </header>

        {articles.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-2xl text-gray-500">記事が見つかりませんでした</p>
          </div>
        ) : (
          <div className="space-y-6">
            {articles.map((article) => (
              <div 
                key={article.id} 
                className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300 overflow-hidden"
              >
                <div className="p-6">
                  <h2 className="text-2xl font-bold text-gray-900 mb-3">
                    <a 
                      href={article.link} 
                      target="_blank" 
                      rel="noopener noreferrer" 
                      className="text-blue-600 hover:text-blue-800 transition-colors"
                    >
                      {article.title}
                    </a>
                  </h2>
                  
                  {article.summary && (
                    <p className="text-gray-600 mb-4 line-clamp-3">
                      {article.summary}
                    </p>
                  )}

                  <div className="flex items-center justify-between text-sm text-gray-500">
                    <div className="flex items-center space-x-3">
                      <div className="flex items-center">
                        <FaClock className="mr-2 text-gray-400" />
                        <span>
                          {new Date(article.published_at).toLocaleString()}
                        </span>
                      </div>
                    </div>
                    <a 
                      href={article.link} 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-blue-500 hover:text-blue-700 flex items-center"
                    >
                      記事を読む
                      <svg 
                        className="w-4 h-4 ml-2" 
                        fill="none" 
                        stroke="currentColor" 
                        viewBox="0 0 24 24" 
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path 
                          strokeLinecap="round" 
                          strokeLinejoin="round" 
                          strokeWidth={2} 
                          d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" 
                        />
                      </svg>
                    </a>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
