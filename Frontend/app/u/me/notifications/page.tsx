"use client"
import Loading from '@/app/components/loading';
import NotificationsAPI from '@/app/utils/apis/notifications'
import { TNotificationPrefs } from '@/app/utils/types';
import React, { useEffect, useState } from 'react'
import { toast } from 'react-toastify';

const Notifications = () => {
    const { data } = NotificationsAPI.getNotificationsPrefs();
    const [prefs, setPrefs] = useState<TNotificationPrefs>();
    useEffect(() => {
        if (data) {
            setPrefs(data.data);
        }
    }, [data]);

    const { mutateAsync } = NotificationsAPI.updateUserPrefs();

    const handleClickUpdatePrefs = async (data: any) => {
        await toast.promise(mutateAsync(data), {
            pending: "Updating Prefs",
            error: "Something went wrong...",
            success: "Prefs updated!"
        }).then((res) => {
            console.log(res.data)
            setPrefs({ ...prefs, ...data })
        })
    }
    return <section>
        <div className='flex flex-col justify-start'>
            <h3 className='text-2xl text-white leading-relaxed'>Notification settings</h3>
            <div className='pt-2 pb-4'>
                <p className='text-sm text-white font-medium'>Choose which types of notifications you want to receive via mail and in-app notifications</p>
            </div>
        </div>

        {
            prefs ? <div className='w-1/2 py-6'>
                <div className="flex items-center p-2 px-4 border border-border rounded-sm bg-transparent">
                    <label className="w-full py-2 text-sm font-bold text-white/80">Asset Sales</label>
                    <input onChange={(e) => handleClickUpdatePrefs({ sale: e.target.checked })} checked={prefs.sale} type="checkbox" className="w-5 h-5 text-white bg-secondary border-secondary rounded-sm" />
                </div>
                <div className="flex items-center p-2 px-4 border border-border rounded-sm bg-transparent">
                    <label className="w-full py-2 text-sm font-bold text-white/80">Asset Verification</label>
                    <input onChange={(e) => handleClickUpdatePrefs({ verified: e.target.checked })} checked={prefs.verified} type="checkbox" className="w-5 h-5 text-white bg-secondary border-secondary rounded-sm" />
                </div>
                <div className="flex items-center p-2 px-4 border border-border rounded-sm bg-transparent">
                    <label className="w-full py-2 text-sm font-bold text-white/80">Asset Rejection</label>
                    <input onChange={(e) => handleClickUpdatePrefs({ rejected: e.target.checked })} checked={prefs.rejected} type="checkbox" className="w-5 h-5 text-white bg-secondary border-secondary rounded-sm" />
                </div>
                <div className="flex items-center p-2 px-4 border border-border rounded-sm bg-transparent">
                    <label className="w-full py-2 text-sm font-bold text-white/80">Login</label>
                    <input onChange={(e) => handleClickUpdatePrefs({ login: e.target.checked })} checked={prefs.login} type="checkbox" className="w-5 h-5 text-white bg-secondary border-secondary rounded-sm" />
                </div>

            </div> : <Loading />
        }

    </section>
}

export default Notifications