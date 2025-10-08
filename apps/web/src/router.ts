// Generouted, changes to this file will be overridden
/* eslint-disable */

import { components, hooks, utils } from '@generouted/react-router/client'

export type Path =
  | `/`
  | `/host/events`
  | `/host/events/:eventId`
  | `/host/events/:eventId/edit`
  | `/host/events/create`
  | `/host/home`
  | `/oauth/google`
  | `/onboard/:method`
  | `/signin`
  | `/signup`

export type Params = {
  '/host/events/:eventId': { eventId: string }
  '/host/events/:eventId/edit': { eventId: string }
  '/onboard/:method': { method: string }
}

export type ModalPath = never

export const { Link, Navigate } = components<Path, Params>()
export const { useModals, useNavigate, useParams } = hooks<Path, Params, ModalPath>()
export const { redirect } = utils<Path, Params>()
