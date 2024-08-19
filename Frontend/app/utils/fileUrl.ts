const getFileUrl = (file: File) => {
  return URL.createObjectURL(file);
};

export default getFileUrl;
