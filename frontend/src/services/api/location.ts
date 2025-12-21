import { Get } from './https';
import type { Location } from '@/interfaces/interview';

export const LocationAPI = {
    getAllLocations: async (): Promise<Location[]> => {
        const response = await Get('/locations');
        return response;
    },
};
