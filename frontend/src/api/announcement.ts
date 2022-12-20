import announcements from '../mocks/announcements.json'
import { IAnnouncement } from "../interfaces/announcement";
import { api } from "./config";
import { errorMessages } from '../constants/messages';

export function getAnnouncements(sku: string, token: string): Promise<IAnnouncement[]> {
  const header = {
    'Authorization': `Bearer ${token}`
  }

  return new Promise((res, rej) => {
    api
      .get(`/announcement/${sku}`, { headers: header })
      .then((resp) =>
        res(
          resp.data.map((announcement: IAnnouncement): IAnnouncement => ({
            id: announcement.id,
            picture: announcement.picture,
            title: announcement.title,
            price: announcement.price,
            quantity: announcement.quantity,
            sku: announcement.sku,
            link: announcement.link
          }))
        )
      )
      .catch((err) => {
        if (err.status === 404) {
          rej({
            ...err,
            message: errorMessages.ANNOUNCEMENT_NOT_FOUND,
          });
        } else {
          rej(err);
        }
      });
  })
}

export function clone(rootID: string, titles: string[], token: string): Promise<void> {
  const body = {
    root_id: rootID,
    titles: titles,
  }

  const header = {
    'Authorization': `Bearer ${token}`
  }
  return new Promise((res, rej) => {
    api
      .post(`/announcement`, body, { headers: header })
      .then(_ => res())
      .catch((err) => {
        if (err.status === 404) {
          rej({
            ...err,
            message: errorMessages.ANNOUNCEMENT_NOT_FOUND,
          });
        } else {
          rej(err);
        }
      });
  })
}