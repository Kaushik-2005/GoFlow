package main

func startJob(store JobReaderWriter, id string) bool {
	job, ok := store.Get(id)
	if !ok {
		return false
	}

	job.MarkRunning()

	return store.Update(job)
}
