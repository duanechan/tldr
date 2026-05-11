import { Textarea } from "@/components/ui/textarea";

export default function Home() {
  return (
    <div className="flex flex-col justify-between items-center h-svh p-8">
      <h1 className="text-4xl font-playfair font-bold">
        Too Long, Didn't Read
      </h1>
      <div className="flex flex-col w-full gap-2">
        <Textarea
          id="contentTextField"
          rows={10}
          placeholder="Paste some contents to summarize"
        />
      </div>
    </div>
  );
}
