import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupTextarea,
} from "@/components/ui/input-group";
import {
  FileIcon,
  FileJpgIcon,
  FilePdfIcon,
  FilePngIcon,
  FileTextIcon,
  GifIcon,
  MarkdownLogoIcon,
  PlusIcon,
  type Icon,
} from "@phosphor-icons/react";
import { useState } from "react";

export default function Home() {
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);

  function handleFileUpload() {
    const input = document.createElement("input");
    input.type = "file";
    input.style.display = "none";
    document.body.appendChild(input);
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (!file) return;
      setUploadedFile(file);
      document.body.removeChild(input);
    };
    input.click();
  }

  return (
    <div className="flex flex-col justify-between items-center h-svh p-8">
      <h1 className="text-4xl font-playfair font-bold">
        Too Long, Didn't Read
      </h1>
      <div className="flex flex-col w-full gap-2">
        <InputGroup>
          <InputGroupTextarea
            id="contentTextField"
            rows={10}
            placeholder="Paste some content to summarize..."
            className="px-6 py-4"
          />
          <InputGroupAddon align="block-end" onClick={handleFileUpload}>
            <InputGroupButton>
              {uploadedFile ? (
                <>
                  <FileTypeIcon fileType={uploadedFile.type} />
                  <span>
                    <span>{`${uploadedFile.name.slice(0, 8)}...${uploadedFile.name.split(".").pop()}`}</span>
                  </span>
                </>
              ) : (
                <>
                  {}
                  <PlusIcon className="cursor-pointer" />
                  <span>Upload a file</span>
                </>
              )}
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </div>
    </div>
  );
}

function FileTypeIcon({ fileType }: { fileType: string }) {
  const Icon = renderFileIcon(fileType);
  return <Icon className="cursor-pointer" />;
}

function renderFileIcon(fileType: string): Icon {
  switch (fileType) {
    case "image/png":
      return FilePngIcon;
    case "image/jpeg":
      return FileJpgIcon;
    case "image/gif":
      return GifIcon;
    case "text/plain":
      return FileTextIcon;
    case "text/markdown":
      return MarkdownLogoIcon;
    case "application/pdf":
      return FilePdfIcon;
    default:
      return FileIcon;
  }
}
