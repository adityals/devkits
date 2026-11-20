export type Lang = "json" | "go" | "ts";

export type Route = `${Lang}_${Lang}`;

export type ConvertOptions = {
  typeName: string;
  goPackage: string;
};

export interface Provider {
  convert(input: string, from: Lang, to: Lang, opts: ConvertOptions): Promise<string>;
}
