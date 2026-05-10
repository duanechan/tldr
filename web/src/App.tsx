import type React from "react";
import { useState } from "react";
import ReactMarkdown from "react-markdown";
import { DragDropInput } from "./components/DragDropInput";

function App() {
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
  const [response, setResponse] = useState("");
  const [loading, setLoading] = useState(false);

  async function summarizeDocument(e: React.SubmitEvent) {
    e.preventDefault();
    if (!uploadedFile) return;

    try {
      setLoading(true);

      const form = new FormData();
      form.append("document", uploadedFile);

      const res = await fetch(
        "http://localhost:8080/api/v1/summarize/document",
        {
          method: "POST",
          body: form,
        },
      );
      const { response } = await res.json();
      setResponse(response);
      console.log(response);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="px-32 py-8">
      <section className="">
        <form className="flex flex-col" onSubmit={summarizeDocument}>
          <DragDropInput
            onChange={(files) => setUploadedFile(files?.item(0) ?? null)}
          />
          <button
            className="px-4 py-2 bg-emerald-300 font-semibold rounded"
            type="submit"
          >
            Summarize
          </button>
        </form>
      </section>
      {loading ? (
        <span>Summarizing...</span>
      ) : (
        response && (
          <section>
            <h1 className="text-2xl font-bold py-8">Summary</h1>
            <div className="prose prose-invert max-w-none border border-gray-600 rounded-xl p-4 text-justify">
              <ReactMarkdown>{response}</ReactMarkdown>
            </div>
          </section>
        )
      )}
    </main>
  );
}

export default App;
