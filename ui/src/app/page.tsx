import CaddyEditor from "@/app/CaddyEditor";
import {Suspense} from "react";
import HistoryViewer from "@/app/HistoryViewer";
import LogStream from "@/app/LogStream";

export default function Home() {
  return (
    <div className="w-full min-h-[100svh]">
      {/*<HistoryViewer/>*/}
      {/*<br/>*/}
      {/*<CaddyEditor />*/}
      <LogStream />
    </div>
  );
}
