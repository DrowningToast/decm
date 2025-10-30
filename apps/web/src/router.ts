// Generouted, changes to this file will be overridden
 

import { components, hooks, utils } from "@generouted/react-router/client";

export type Path =
    | `/`
    | `/app`
    | `/auth/success`
    | `/error`
    | `/host/events`
    | `/host/events/:eventId`
    | `/host/events/:eventId/edit`
    | `/host/events/:eventId/settings/certificate`
    | `/host/events/:eventId/settings/participant`
    | `/host/events/create`
    | `/host/home`
    | `/onboard/:method`
    | `/signin`
    | `/signout`
    | `/signup`;

export type Params = {
    "/host/events/:eventId": { eventId: string };
    "/host/events/:eventId/edit": { eventId: string };
    "/host/events/:eventId/settings/certificate": { eventId: string };
    "/host/events/:eventId/settings/participant": { eventId: string };
    "/onboard/:method": { method: string };
};

export type ModalPath = never;

export const { Link, Navigate } = components<Path, Params>();
export const { useModals, useNavigate, useParams } = hooks<Path, Params, ModalPath>();
export const { redirect } = utils<Path, Params>();
