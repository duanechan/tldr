import { marked } from "marked";
import type React from "react";
import { useState } from "react";

function App() {
  const [response, setResponse] = useState("");
  const [loading, setLoading] = useState(false);

  async function summarizeDocument(e: React.SubmitEvent) {
    e.preventDefault();

    const file = e.target.elements.namedItem("document") as HTMLInputElement;
    const uploadedFile = file.files?.[0];
    if (!uploadedFile) return;

    try {
      setLoading(true);

      const form = new FormData();
      form.append("document", uploadedFile);

      const res = await fetch("/api/v1/summarize/document", {
        method: "POST",
        body: form,
      });
      const { response } = await res.json();
      setResponse(await marked.parse(response));
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="px-32 py-8">
      <section>
        <form className="flex flex-col w-50" onSubmit={summarizeDocument}>
          <label htmlFor="documentField">Document</label>
          <input id="documentField" name="document" type="file" />
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
        <section>
          <h1>Summary</h1>
          <div dangerouslySetInnerHTML={{ __html: response }}></div>
        </section>
      )}
    </main>
  );
}

export default App;
