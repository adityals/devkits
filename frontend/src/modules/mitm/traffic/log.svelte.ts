import { Events } from "@wailsio/runtime";
import {
  Capacity,
  Clear,
  List,
  Recording,
  SetRecording,
} from "../../../../bindings/devkits/internal/mitm/traffic/trafficservice";
import type { Summary } from "../../../../bindings/devkits/internal/mitm/traffic";
import { MitmTrafficFlow } from "../../../../bindings/events";

export type Row = Summary & {
  host: string;
  path: string;
};

const toRow = (s: Summary): Row => {
  const url = URL.parse(s.url);

  return {
    ...s,
    host: url?.host ?? "",
    path: url ? url.pathname + url.search : s.url,
  };
};

export class TrafficLog {
  rows = $state.raw<Row[]>([]);
  error = $state("");
  recording = $state(true);
  #maxRows = Infinity;

  start(): () => void {
    const off = Events.On(MitmTrafficFlow, (event) => {
      this.rows = [...this.rows, toRow(event.data)].slice(-this.#maxRows);
    });

    Promise.all([Capacity(), Recording(), List()])
      .then(([capacity, recording, list]) => {
        this.#maxRows = capacity;
        this.recording = recording;
        const live = new Set(this.rows.map((r) => r.id));
        const past = (list ?? []).filter((s) => !live.has(s.id)).map(toRow);

        this.rows = [...past, ...this.rows].slice(-this.#maxRows);
      })
      .catch(this.#fail);

    return off;
  }

  toggleRecording(): void {
    const next = !this.recording;

    SetRecording(next)
      .then(() => {
        this.recording = next;
      })
      .catch(this.#fail);
  }

  clear(): void {
    Clear()
      .then(() => {
        this.rows = [];
      })
      .catch(this.#fail);
  }

  #fail = (error: unknown): void => {
    this.error = String(error);
  };
}
