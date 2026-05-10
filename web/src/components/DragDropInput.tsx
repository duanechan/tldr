import { useRef, useState } from "react";

type DragDropInputProps = {
  onChange: (files: FileList | null) => void;
};

export function DragDropInput({ onChange }: DragDropInputProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [files, setFiles] = useState<FileList | null>(null);
  const [isDragOver, setDragOver] = useState(false);

  return (
    <div
      onDragOver={(e) => {
        e.preventDefault();
        setDragOver(true);
      }}
      onDrop={(e) => {
        e.preventDefault();
        const dropped = e.dataTransfer.files;
        setFiles(dropped);
        onChange(dropped);
        setDragOver(false);
      }}
      onDragEnter={() => setDragOver(true)}
      onDragLeave={() => setDragOver(false)}
      onClick={() => inputRef.current?.click()}
      className={`${isDragOver ? "border-2 border-blue-400" : "border-2 border-transparent"} 
        flex flex-col items-center justify-center h-[50vh] bg-gray-900
        rounded-2xl
      `}
    >
      <input
        type="file"
        className="hidden"
        ref={inputRef}
        onChange={(e) => {
          e.preventDefault();
          if (!e.target.files) return;
          setFiles(e.target.files);
          onChange(e.target.files);
        }}
      />
      {files
        ? files.length > 1
          ? `${files.item(0)?.name}, and ${files.length - 1} others...`
          : `${files.item(0)?.name}`
        : "Drop files here"}
    </div>
  );
}
