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
  capacity = $state(0);

  start(): () => void {
    const off = Events.On(MitmTrafficFlow, (event) => {
      this.rows = this.#trim([...this.rows, toRow(event.data)]);
    });

    Promise.all([Capacity(), Recording(), List()])
      .then(([capacity, recording, list]) => {
        this.capacity = capacity;
        this.recording = recording;
        const live = new Set(this.rows.map((r) => r.id));
        const past = (list ?? []).filter((s) => !live.has(s.id)).map(toRow);

        this.rows = this.#trim([...past, ...this.rows]);
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

  // Before Capacity() resolves, capacity is 0 and rows are kept untrimmed.
  #trim = (rows: Row[]): Row[] =>
    this.capacity > 0 ? rows.slice(-this.capacity) : rows;

  #fail = (error: unknown): void => {
    this.error = String(error);
  };
}
