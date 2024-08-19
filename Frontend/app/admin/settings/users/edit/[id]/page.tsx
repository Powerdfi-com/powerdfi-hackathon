"use client"
import Loading from '@/app/components/loading'
import AdminAPI from '@/app/utils/apis/admin'
import { useRouter } from 'next/navigation'
import React, { useEffect, useState } from 'react'
import { GoPerson } from 'react-icons/go'
import { toast } from 'react-toastify'

const EditAdmin = ({ params }: { params: { id: string } }) => {
    const router = useRouter();
    const [data, setData] = useState<{ email: string, password: string; confirmPassword: string; role?: number }>({
        email: "",
        password: "",
        confirmPassword: "",
    })

    const { data: roles } = AdminAPI.getRoles()
    const { data: admin } = AdminAPI.getAdminById(params.id);
    const { mutateAsync, isPending } = AdminAPI.updateAdmin()
    const handleClickUpdateAdmin = async (e: any) => {
        e.preventDefault()
        if (!data.email || !data.role || !data.password || !data.confirmPassword) {
            toast.error("Please fill in details properly!");
            return
        }
        const { email, password, role } = data;
        await toast.promise(mutateAsync({ email, password, role, id: params.id }), {
            pending: "Updating admin, please wait!",
            error: "Error updating admin!",
            success: "Admin updated successfully!"
        }).then((res) => {
            router.back()
        })
    }
    useEffect(() => {
        if (admin && roles) {
            setData((_data) => {
                return { ..._data, email: admin.data.email, role: roles.data.find((roles) => roles.name === admin.data.roles[0])!.id }
            })
        }
    }, [admin, roles])
    return <section>
        <div className='flex justify-between'>
            <div className='text-white flex !text-2xl gap-4'><GoPerson /> <span>Edit User</span></div>
        </div>
        <h4 className='text-white mt-8'>User Account Info</h4>
        {
            (roles && admin) ? <form className="flex-1 flex flex-col gap-2">
                <label className='flex flex-col gap-2 mt-3'>
                    <span className='text-sm text-white'>Email</span>
                    <input value={data.email} onChange={(e) => setData({ ...data, email: e.target.value })} type='email' className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' />
                </label>
                <label className='flex flex-col gap-2 mt-3'>
                    <span className='text-sm text-white'>Password</span>
                    <input value={data.password} onChange={(e) => setData({ ...data, password: e.target.value })} type='password' className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' />
                </label>
                <label className='flex flex-col gap-2 mt-3'>
                    <span className='text-sm text-white'>Confirm Password</span>
                    <input value={data.confirmPassword} onChange={(e) => setData({ ...data, confirmPassword: e.target.value })} type='password' className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' />
                </label>
                <label className='flex flex-col gap-2 mt-3'>
                    <span className='text-sm text-white'>Role Access</span>
                    <select onChange={(e) => setData({ ...data, role: parseFloat(e.target.value) })} className='max-w-sm bg-transparent border-text-grey border rounded-md text-white py-2 text-sm px-5' >
                        {roles?.data.map((role) => <option value={role.id} key={role.id}>{role.name}</option>)}
                    </select>
                </label>
                <div className="flex gap-6 w-full mt-6 mb-8 max-w-xs">
                    <button disabled={isPending} onClick={handleClickUpdateAdmin} className="text-sm bg-secondary py-3 rounded-md flex-1">Save</button>
                    <button disabled={isPending} className="text-sm bg-transparent ring-1 ring-primary/40 text-secondary py-3 rounded-md flex-1" onClick={() => router.back()}>Cancel</button>
                </div>
            </form> : <Loading />
        }
    </section>
}

export default EditAdmin