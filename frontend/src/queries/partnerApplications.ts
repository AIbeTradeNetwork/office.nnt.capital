import { $request } from 'utils/axios'

export interface PartnerApplication {
  uid: string
  applicant_uid: string
  partner_uid: string
  status: 'pending' | 'approved' | 'rejected'
  message: string
  response: string
  created_at: number
  processed_at?: number
  processed_by: string
  applicant?: {
    uid: string
    nickname: string
    email: string
  }
  partner?: {
    uid: string
    nickname: string
    email: string
  }
}

export interface PartnerApplicationReq {
  partner_uid: string
  message: string
}

export interface PartnerApplicationResponseReq {
  application_uid: string
  status: 'pending' | 'approved' | 'rejected'
  response: string
}

// Получить заявки на партнёрство (для партнёра)
export async function getPartnerApplications(limit: number, skip: number): Promise<PartnerApplication[]> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `query ($limit: Int!, $skip: Int!) {
        partner_applications(limit: $limit, skip: $skip) {
          uid
          applicant_uid
          partner_uid
          status
          message
          response
          created_at
          processed_at
          processed_by
          applicant {
            uid
            nickname
            email
          }
          partner {
            uid
            nickname
            email
          }
        }
      }`,
      variables: {
        limit,
        skip,
      },
      cache: false,
    })

    return response.data.data?.partner_applications || []
  } catch (error) {
    return Promise.reject(error)
  }
}

// Получить мои заявки на партнёрство
export async function getMyApplications(limit: number, skip: number): Promise<PartnerApplication[]> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `query ($limit: Int!, $skip: Int!) {
        my_applications(limit: $limit, skip: $skip) {
          uid
          applicant_uid
          partner_uid
          status
          message
          response
          created_at
          processed_at
          processed_by
          applicant {
            uid
            nickname
            email
          }
          partner {
            uid
            nickname
            email
          }
        }
      }`,
      variables: {
        limit,
        skip,
      },
      cache: false,
    })

    return response.data.data?.my_applications || []
  } catch (error) {
    return Promise.reject(error)
  }
}

// Получить количество заявок на партнёрство (для партнёра)
export async function getPartnerApplicationsCount(): Promise<number> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `query {
        partner_applications_count
      }`,
      cache: false,
    })

    return response.data.data?.partner_applications_count || 0
  } catch (error) {
    return Promise.reject(error)
  }
}

// Получить количество моих заявок на партнёрство
export async function getMyApplicationsCount(): Promise<number> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `query {
        my_applications_count
      }`,
      cache: false,
    })

    return response.data.data?.my_applications_count || 0
  } catch (error) {
    return Promise.reject(error)
  }
}

// Создать заявку на партнёрство
export async function createPartnerApplication(input: PartnerApplicationReq): Promise<PartnerApplication> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `mutation ($input: PartnerApplicationReq!) {
        create_partner_application(input: $input) {
          uid
          applicant_uid
          partner_uid
          status
          message
          response
          created_at
          processed_at
          processed_by
          applicant {
            uid
            nickname
            email
          }
          partner {
            uid
            nickname
            email
          }
        }
      }`,
      variables: {
        input,
      },
      cache: false,
    })

    return response.data.data?.create_partner_application
  } catch (error) {
    return Promise.reject(error)
  }
}

// Обработать заявку на партнёрство (одобрить или отклонить)
export async function processPartnerApplication(input: PartnerApplicationResponseReq): Promise<PartnerApplication> {
  try {
    const response = await $request({
      requestConfig: {
        url: import.meta.env.VITE_APP_URL_GQL,
      },
      query: `mutation ($input: PartnerApplicationResponseReq!) {
        process_partner_application(input: $input) {
          uid
          applicant_uid
          partner_uid
          status
          message
          response
          created_at
          processed_at
          processed_by
          applicant {
            uid
            nickname
            email
          }
          partner {
            uid
            nickname
            email
          }
        }
      }`,
      variables: {
        input,
      },
      cache: false,
    })

    return response.data.data?.process_partner_application
  } catch (error) {
    return Promise.reject(error)
  }
} 