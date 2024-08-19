"use client"
import FileUpload from "@/app/components/fileUpload";
import Loading from "@/app/components/loading";
import AssetAPI from "@/app/utils/apis/asset";
import { useState } from "react";
import Error from '@/app/components/error';
import { toast } from "react-toastify";
import { uploadFile } from "@/app/utils/apis/file-upload";
import { useRouter } from "next/navigation";
import Image from "next/image";
import getFileUrl from '@/app/utils/fileUrl';
import { generateArray } from "@/app/utils/func";
import { IoMdAdd } from "react-icons/io";


const CreateAsset = ({ params }: { params: { blockchainId: string } }) => {
    const [files, setFiles] = useState<FileList>()
    const [legalFiles, setLegalFiles] = useState<FileList>()
    const [issuanceFiles, setIssuanceFiles] = useState<FileList>();
    const { isPending, data, isSuccess } = AssetAPI.getCategories();
    const [divisible, setDivivisible] = useState(true);
    const router = useRouter()
    const [asset, setAsset] = useState<{ name: string, symbol: string, totalSupply: string, description: string, categoryId: number, properties: [string, string][], blockchainId: string }>({
        name: "",
        symbol: "",
        totalSupply: "1",
        description: "",
        categoryId: 1,
        properties: [],
        blockchainId: params.blockchainId
    });
    const [property, setProperty] = useState({
        key: "",
        value: ""
    })
    const { mutateAsync: createAsset, isPending: isCreatingAsset } = AssetAPI.create({
        onSuccess: () => {
            router.push("/u/me/portfolio")
        }
    })
    const handleClickCreateToken = (e: any) => {
        e.preventDefault();
        const createAssetPromise = new Promise(async (resolve, reject) => {
            try {
                const urls = await uploadFile(files!);
                const legalDocs = await uploadFile(legalFiles!);
                const issuanceDocs = await uploadFile(issuanceFiles!);
                await createAsset({ ...asset, totalSupply: parseInt(asset.totalSupply), properties: JSON.stringify(asset.properties), urls, legalDocs, issuanceDocs });
                resolve("")
            } catch (e) {
                reject()
            }
        })
        toast.promise(createAssetPromise, {
            error: "Asset creation failed!",
            pending: "Creating asset, please wait!",
            success: "Asset successfully created!"
        })
    }

    const handleClickAddProperty = (e: any) => {
        e.preventDefault()
        if (property.key && property.value) {
            if (asset.properties.find((ass) => ass[0] === property.key)) {
                toast.error("Property with key already exists!");
                return;
            } setAsset({ ...asset, properties: [...asset.properties, [property.key, property.value]] })
            setProperty({ key: "", value: "" });
        }
    }

    const removeProperty = (key: string) => {
        setAsset({ ...asset, properties: asset.properties.filter((ass) => ass[0] !== key) })
    }
    return (
        <section className="">
            <h3 className="text-white text-2xl leading-relaxed">Create Asset Token</h3>
            {
                isPending ? <Loading /> : (data ? <section className="flex mt-8">
                    <section className="flex-1">
                        <h4 className="text-white text-lg leading-relaxed">Upload Asset Image</h4>
                        <p className="text-sm text-text-grey">File types supported: JPG and PNG Max 100mb.</p>
                        {
                            files && <FileDisplay files={files} />
                        }
                        <FileUpload multiple={true} updateValue={(e) => setFiles(e)}>
                            <div className="h-48 w-full max-w-sm ring-1 ring-black-shade flex flex-col items-center justify-center rounded-xl mt-6">
                                <span className="px-3 py-1 rounded-full ring-1 ring-primary/40 text-white text-sm">Choose files</span>
                            </div>
                        </FileUpload>
                    </section>
                    <form className="flex-1 flex flex-col gap-2">
                        <label className='flex flex-col gap-2 mt-3'>
                            <span className='text-sm text-white'>Asset Name</span>
                            <input className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' placeholder='Elton Ave' value={asset.name} onChange={(e) => setAsset({ ...asset, name: e.target.value })} />
                        </label>
                        <label className='flex flex-col gap-2 mt-3'>
                            <span className='text-sm text-white'>Asset Symbol</span>
                            <input className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' placeholder='Elton Ave' value={asset.symbol} onChange={(e) => setAsset({ ...asset, symbol: e.target.value })} />
                        </label>
                        <div className='flex flex-col mt-3 gap-2'>
                            <h4 className='text-sm text-white'>Submit legal Asset document</h4>
                            <p className="text-xs text-text-grey">File types supported: PDF 5.5mb.</p>
                            {
                                legalFiles && <FileDisplay files={legalFiles} />
                            }
                            <FileUpload multiple={true} updateValue={(e) => setLegalFiles(e)}>
                                <div className="h-16 max-w-sm w-full ring-1 ring-black-shade flex flex-col items-center justify-center rounded-xl">
                                    <span className="px-3 py-1 rounded-full ring-1 ring-primary/40 text-white text-sm">Choose files</span>
                                </div>
                            </FileUpload>
                        </div>
                        <div className='flex flex-col mt-3 gap-2'>
                            <h4 className='text-sm text-white'>Submit your digital asset issuance documents</h4>
                            <p className="text-xs text-text-grey">File types supported: PDF 5.5mb.</p>
                            {
                                issuanceFiles && <FileDisplay files={issuanceFiles} />
                            }
                            <FileUpload multiple={true} updateValue={(e) => setIssuanceFiles(e)}>
                                <div className="h-16 max-w-sm w-full ring-1 ring-black-shade flex flex-col items-center justify-center rounded-xl">
                                    <span className="px-3 py-1 rounded-full ring-1 ring-primary/40 text-white text-sm">Choose files</span>
                                </div>
                            </FileUpload>
                        </div>
                        <label className='flex flex-col gap-2 mt-3'>
                            <span className='text-sm text-white'>Description</span>
                            <textarea className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' placeholder='Description goes here' value={asset.description} onChange={(e) => setAsset({ ...asset, description: e.target.value })} />
                        </label>
                        <label className='flex flex-col gap-2 mt-3'>
                            <span className='text-sm text-white'>Asset Type</span>
                            <select className='max-w-sm bg-transparent border-text-grey border rounded-md text-white py-2 text-sm px-5' value={asset.categoryId} onChange={(e) => setAsset({ ...asset, categoryId: parseInt(e.target.value) })} >
                                {data.data.map((category) => <option key={category.id} value={category.id}>{category.name.toUpperCase()}</option>)}
                            </select>
                        </label>
                        <label className='flex flex-col gap-2 mt-3  max-w-sm'>
                            <span className='text-sm text-white'>Asset Features</span>
                            <ul className="flex flex-col gap-2">
                                {asset.properties.map((ass) => <li key={ass[0]}>
                                    <div className="flex gap-2 items-center w-full">
                                        <input className='indent-2 flex-[2] text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm' placeholder='Key' value={ass[0]} disabled />
                                        <input className='indent-2 flex-[2] text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm' placeholder='Value' value={ass[1]} disabled />
                                        <button className="flex-[1]  border rounded-md border-text-grey bg-transparent py-2 text-xs text-white/40" onClick={() => removeProperty(ass[0])}>Remove</button>
                                    </div>
                                </li>)}
                            </ul>
                            <div className="flex gap-2 items-center w-full">
                                <input className='indent-2 flex-[2] text-white outline-none border rounded-md border-white/60 bg-transparent py-2 text-sm' placeholder='Key' value={property.key} onChange={(e) => setProperty({ ...property, key: e.target.value })} />
                                <input className='indent-2 flex-[2] text-white outline-none border rounded-md border-white/60 bg-transparent py-2 text-sm' placeholder='Value' value={property.value} onChange={(e) => setProperty({ ...property, value: e.target.value })} />
                                <button className="flex-[1] flex justify-center items-center  border rounded-md border-white/60 bg-transparent py-2 text-sm" onClick={handleClickAddProperty}><IoMdAdd className="text-white !text-xl" /></button>
                            </div>
                        </label>
                        <div className='flex flex-col gap-2  mt-3'>
                            <h4 className='text-sm text-white'>Asset Divisibility</h4>
                            <div>
                                <label className="flex gap-2 items-center"><input type="radio" name="divisibility" className='text-sm text-white' value={"yes"} onChange={(e) => setDivivisible(true)} checked={divisible} /><span className="text-sm text-text-grey">Yes</span></label>
                                <label className="flex gap-2 items-center"><input type="radio" name="divisibility" className='text-sm text-white' value={"no"} onChange={(e) => {
                                    setDivivisible(false);
                                    setAsset({ ...asset, totalSupply: "1" })
                                }} checked={!divisible} /><span className="text-sm text-text-grey">No</span></label>
                            </div>
                        </div>
                        {
                            divisible && <label className='flex flex-col gap-2 mt-3'>
                                <span className='text-sm text-white'>Total Supply</span>
                                <input type="number" className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' placeholder='100' value={asset.totalSupply} onChange={(e) => setAsset({ ...asset, totalSupply: e.target.value })} />
                            </label>
                        }
                        <div className="flex gap-2 w-full mt-6 mb-8 max-w-sm">
                            <button disabled={isCreatingAsset} className="text-sm bg-transparent ring-1 ring-primary/40 text-secondary h-8 rounded-md flex-1" onClick={() => router.back()}>Cancel</button>
                            <button className="text-sm bg-secondary h-8 rounded-md flex-1" onClick={handleClickCreateToken}>Create Token</button>
                        </div>
                    </form>
                </section> : <Error />)
            }
        </section>
    )
}

export default CreateAsset;

const FileDisplay = ({ files }: { files: FileList }) => {
    return <ul className="flex gap-2 w-full my-2 flex-wrap">
        {
            generateArray(files.length).map((_, index) => {
                const file = files!.item(index)!;
                const url = getFileUrl(file);
                return <li key={file.name}><div className="h-10 w-10 rounded-md relative"><Image src={url} alt="" fill={true} className="rounded-md object-cover" /></div></li>
            })
        }
    </ul>
}