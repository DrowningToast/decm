// Generouted, changes to this file will be overridden
/* eslint-disable */

import { components, hooks, utils } from "@generouted/react-router/client";

export type Path =
    | `/`
    | `/app`
    | `/app/certificates`
    | `/app/certificates/:id`
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
    | `/signin/sign-message`
    | `/signin/verify-oauth`
    | `/signout`
    | `/signup`;

export type Params = {
    "/app/certificates/:id": { id: string };
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
