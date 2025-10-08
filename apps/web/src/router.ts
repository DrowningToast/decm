// Generouted, changes to this file will be overridden
/* eslint-disable */

import { components, hooks, utils } from '@generouted/react-router/client'

export type Path =
  | `/`
  | `/host/events`
  | `/host/events/create`
  | `/host/home`
  | `/oauth/google`
  | `/onboard/:method`
  | `/signin`
  | `/signup`

export type Params = {
  '/onboard/:method': { method: string }
}

export type ModalPath = never

export const { Link, Navigate } = components<Path, Params>()
export const { useModals, useNavigate, useParams } = hooks<Path, Params, ModalPath>()
export const { redirect } = utils<Path, Params>()
