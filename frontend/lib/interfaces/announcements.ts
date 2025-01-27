export enum Status {
  Ativo = "active",
  Inativo = "closed",
  Pausado = "paused",
}

export interface IAnnouncement {
  id: string,
  title: string;
  quantity: number;
  status: Status;
  price: number;
  picture: string;
  sku: string;
  link: string;
  account: {
    id: string;
    name: string;
  };
}